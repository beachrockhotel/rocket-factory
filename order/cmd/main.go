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
		// ← ТОЖЕ КАК РЕЗУЛЬТАТ, error=nil
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

// Заглушки, чтобы сервис запускался

func (h *OrderHandler) CreatePayment(ctx context.Context, req *order_v1.PayOrderRequest, p order_v1.CreatePaymentParams) (order_v1.CreatePaymentRes, error) {
	return &order_v1.NotFoundError{Code: http.StatusNotFound, Message: "not implemented yet"}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, p order_v1.GetOrderParams) (order_v1.GetOrderRes, error) {
	return &order_v1.NotFoundError{Code: http.StatusNotFound, Message: "not implemented yet"}, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, p order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	return &order_v1.NotFoundError{Code: http.StatusNotFound, Message: "not implemented yet"}, nil
}

// для неожиданных ошибок (когда handler вернул error != nil)
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
