package main

import (
	"context"
	"errors"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	order_v1 "github.com/beachrockhotel/rocket-factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second

	inventoryAddr = "localhost:50051"
	paymentAddr   = "localhost:50052"
)

// ==== In-memory хранилище ====

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*order_v1.OrderDto // key = order_uuid (string)
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{orders: make(map[string]*order_v1.OrderDto)}
}

type OrderHandler struct {
	storage *OrderStorage
	inv     inventory_v1.InventoryServiceClient
	pay     payment_v1.PaymentServiceClient // пока не используется
}

// CreateOrder: POST /api/v1/orders
func (h *OrderHandler) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	// 1) Валидация
	if req.UserUUID == (uuid.UUID{}) {
		return &order_v1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "user_uuid is required",
		}, nil
	}
	if len(req.PartUuids) == 0 {
		return &order_v1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "part_uuids is required",
		}, nil
	}

	// 2) []uuid.UUID -> []string для фильтра
	uuids := make([]string, 0, len(req.PartUuids))
	for _, u := range req.PartUuids {
		uuids = append(uuids, u.String())
	}

	// 3) gRPC в Inventory
	invResp, err := h.inv.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Filter: &inventory_v1.PartsFilter{Uuids: uuids},
	})
	if err != nil {
		return &order_v1.BadGatewayError{
			Code:    http.StatusBadGateway,
			Message: "inventory unavailable",
		}, nil
	}

	// 4) Проверка, что все детали существуют (учёт кратности)
	found := make(map[string]*inventory_v1.Part, len(invResp.GetParts()))
	for _, p := range invResp.GetParts() {
		found[p.GetUuid()] = p
	}
	countBy := make(map[string]int, len(uuids))
	var missing []string
	for _, id := range uuids {
		countBy[id]++
		if _, ok := found[id]; !ok {
			missing = append(missing, id)
		}
	}
	if len(missing) > 0 {
		return &order_v1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "some parts not found: " + strings.Join(missing, ", "),
		}, nil
	}

	// 5) Расчёт total_price (учитываем кратность), округляем до 2 знаков
	var total float64
	for id, n := range countBy {
		total += found[id].GetPrice() * float64(n)
	}
	total = math.Round(total*100) / 100

	// 6) Генерим order_uuid и сохраняем
	orderID := uuid.NewString()
	h.storage.mu.Lock()
	h.storage.orders[orderID] = &order_v1.OrderDto{
		OrderUUID:  orderID,
		UserUUID:   req.UserUUID,
		PartUuids:  append([]uuid.UUID(nil), req.PartUuids...),
		TotalPrice: total,
		Status:     order_v1.OrderStatusPENDINGPAYMENT,
		// TransactionUUID, PaymentMethod — нулевые значения (не установлены) → это ок
	}
	h.storage.mu.Unlock()

	// 7) Ответ (успешный результат)
	return &order_v1.CreateOrderResponse{
		OrderUUID:  uuid.MustParse(orderID),
		TotalPrice: float32(total),
	}, nil
}

func (h *OrderHandler) GetOrder(_ context.Context, params order_v1.GetOrderParams) (order_v1.GetOrderRes, error) {
	id := params.OrderUUID.String()
	h.storage.mu.RLock()
	o, ok := h.storage.orders[id]
	h.storage.mu.RUnlock()
	if !ok {
		return &order_v1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}

	return &order_v1.GetOrderResponse{
		Order: *o,
	}, nil
}

func (h *OrderHandler) CreatePayment(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.CreatePaymentParams) (order_v1.CreatePaymentRes, error) {
	id := params.OrderUUID.String()
	pmOpt := req.PaymentMethod

	// 1) Найти заказ + базовые проверки статуса
	h.storage.mu.RLock()
	o, ok := h.storage.orders[id]
	h.storage.mu.RUnlock()
	if !ok {
		return &order_v1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
	}
	if o.Status == order_v1.OrderStatusCANCELLED {
		// для CreatePayment контракт не допускает 409, поэтому 400
		return &order_v1.BadRequestError{Code: http.StatusBadRequest, Message: "order cancelled"}, nil
	}
	if o.Status == order_v1.OrderStatusPAID {
		return &order_v1.BadRequestError{Code: http.StatusBadRequest, Message: "already paid"}, nil
	}

	// 2) Валидация payment_method
	pm, ok := pmOpt.Get()
	if !ok || pm == order_v1.PaymentMethodUNKNOWN {
		return &order_v1.BadRequestError{Code: http.StatusBadRequest, Message: "invalid payment_method"}, nil
	}

	// 3) Маппинг OpenAPI → gRPC enum
	var grpcPM payment_v1.PaymentMethod
	switch pm {
	case order_v1.PaymentMethodCARD:
		grpcPM = payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case order_v1.PaymentMethodSBP:
		grpcPM = payment_v1.PaymentMethod_PAYMENT_METHOD_SBP
	case order_v1.PaymentMethodCREDITCARD:
		grpcPM = payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case order_v1.PaymentMethodINVESTORMONEY:
		grpcPM = payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return &order_v1.BadRequestError{Code: http.StatusBadRequest, Message: "invalid payment_method"}, nil
	}

	// 4) Вызов платёжки
	pr, err := h.pay.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     id,
		UserUuid:      o.UserUUID.String(),
		PaymentMethod: grpcPM,
	})
	if err != nil {
		// для CreatePayment контракт не допускает 502, поэтому 500
		return &order_v1.InternalServerError{Code: http.StatusInternalServerError, Message: "payment unavailable"}, nil
	}

	// 5) Разбор transaction_uuid из ответа платёжки
	txStr := pr.GetTransactionUuid()
	txID, err := uuid.Parse(txStr)
	if err != nil {
		return &order_v1.InternalServerError{Code: http.StatusInternalServerError, Message: "invalid transaction uuid from payment"}, nil
	}

	// 6) Обновить заказ под Lock (на случай гонки проверяем статус повторно)
	h.storage.mu.Lock()
	if o2, ok := h.storage.orders[id]; ok {
		switch o2.Status {
		case order_v1.OrderStatusPENDINGPAYMENT:
			o2.Status = order_v1.OrderStatusPAID
			o2.TransactionUUID.SetTo(txID) // OptNilUUID
			o2.PaymentMethod.SetTo(pm)     // OptNilPaymentMethod
		case order_v1.OrderStatusPAID:
			// если кто-то уже оплатил — вернуть текущий tx идемпотентно
			if tx, ok := o2.TransactionUUID.Get(); ok {
				h.storage.mu.Unlock()
				return &order_v1.PayOrderResponse{TransactionUUID: tx}, nil
			}
		case order_v1.OrderStatusCANCELLED:
			h.storage.mu.Unlock()
			return &order_v1.BadRequestError{Code: http.StatusBadRequest, Message: "order cancelled"}, nil
		}
	}
	h.storage.mu.Unlock()

	// 7) Возврат успешного ответа с реальным transaction_uuid
	return &order_v1.PayOrderResponse{TransactionUUID: txID}, nil
}

func (h *OrderHandler) CancelOrder(_ context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	id := params.OrderUUID.String()

	h.storage.mu.Lock()
	defer h.storage.mu.Unlock()

	o, ok := h.storage.orders[id]
	if !ok {
		return &order_v1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
	}

	switch o.Status {
	case order_v1.OrderStatusPAID:
		// Оплаченный нельзя отменить 409
		return &order_v1.ConflictError{Code: http.StatusConflict, Message: "order already paid"}, nil

	case order_v1.OrderStatusCANCELLED:
		// Уже отменён — идемпотентно считаем успехом → 204
		return &order_v1.CancelOrderNoContent{}, nil

	case order_v1.OrderStatusPENDINGPAYMENT:
		// Отмена ожидающей оплаты ставим CANCELLED и чистим оплату
		o.Status = order_v1.OrderStatusCANCELLED
		o.TransactionUUID.SetToNull()
		o.PaymentMethod.SetToNull()
		return &order_v1.CancelOrderNoContent{}, nil

	default:
		return &order_v1.ConflictError{Code: http.StatusConflict, Message: "invalid order status"}, nil
	}
}

func (h *OrderHandler) NewError(_ context.Context, err error) *order_v1.GenericErrorStatusCode {
	return &order_v1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: order_v1.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}

func main() {
	// gRPC-клиенты к Inventory/Payment
	invConn, err := grpc.Dial(inventoryAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial inventory %s: %v", inventoryAddr, err)
	}
	defer invConn.Close()
	invClient := inventory_v1.NewInventoryServiceClient(invConn)

	payConn, err := grpc.Dial(paymentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial payment %s: %v", paymentAddr, err)
	}
	defer payConn.Close()
	payClient := payment_v1.NewPaymentServiceClient(payConn)

	storage := NewOrderStorage()
	handler := &OrderHandler{
		storage: storage,
		inv:     invClient,
		pay:     payClient,
	}

	// ogen-сервер (оборачивает handler)
	orderServer, err := order_v1.NewServer(handler)
	if err != nil {
		log.Fatalf("openapi server init error: %v", err)
	}

	// HTTP роутер/сервер
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Mount("/", orderServer)

	srv := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// Запуск
	go func() {
		log.Printf("🚀 HTTP OrderService listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ HTTP server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("✅ Server stopped")
}
