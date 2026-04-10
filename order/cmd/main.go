package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPIV1 "github.com/beachrockhotel/rocket-factory/order/internal/api/orders/v1"
	inventoryClient "github.com/beachrockhotel/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/beachrockhotel/rocket-factory/order/internal/client/grpc/payment/v1"
	orderRepo "github.com/beachrockhotel/rocket-factory/order/internal/repository/order"
	orderSrv "github.com/beachrockhotel/rocket-factory/order/internal/service/order"
	generatedOrderV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/openapi/order/v1"
	generatedInventoryV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/inventory/v1"
	generatedPaymentV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort              = ":8080"
	httpServerReadTimeout = 5 * time.Second

	inventoryServiceAddr = "localhost:50051"
	paymentServiceAddr   = "localhost:50052"
)

func main() {
	inventoryV1Conn, err := grpc.NewClient(
		inventoryServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to inventoryV1 connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryV1Conn.Close(); cerr != nil {
			log.Printf("failed to close inventoryV1Conn: %v", cerr)
		}
	}()

	paymentV1Conn, err := grpc.NewClient(paymentServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to paymentV1 connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentV1Conn.Close(); cerr != nil {
			log.Printf("failed to close paymentV1Conn: %v", cerr)
		}
	}()

	generatedInventoryV1Client := generatedInventoryV1.NewInventoryServiceClient(inventoryV1Conn)
	generatedPaymentV1Client := generatedPaymentV1.NewPaymentServiceClient(paymentV1Conn)

	inventoryV1Client := inventoryClient.NewClient(generatedInventoryV1Client)
	paymentV1Client := paymentClient.NewClient(generatedPaymentV1Client)

	orderRepository := orderRepo.NewRepository()
	orderService := orderSrv.NewService(orderRepository, inventoryV1Client, paymentV1Client)

	r := chi.NewRouter()
	srv := orderAPIV1.NewOrderAPI(orderService)

	handler, err := generatedOrderV1.NewServer(srv)
	if err != nil {
		log.Printf("failed to create ogen handler: %v\n", err)
		return
	}

	r.Mount("/", handler)

	server := &http.Server{
		ReadTimeout: httpServerReadTimeout,
		Addr:        httpPort,
		Handler:     r,
	}

	go func() {
		log.Println("🚀 OrderService HTTP API running on :8080")
		if errServer := server.ListenAndServe(); errServer != nil && !errors.Is(errServer, http.ErrServerClosed) {
			log.Printf("server error: %v\n", errServer)
			return
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("graceful shutdown failed: %v\n", err)
		return
	}

	log.Println("✅ Server stopped")
}
