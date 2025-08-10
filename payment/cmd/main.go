package main

import (
	"context"
	"fmt"
	paymentV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	grpcPort = 50052
)

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func (s *paymentService) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if req.GetOrderUuid() == "" || req.GetUserUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_uuid and user_uuid are required")
	}

	txID := uuid.NewString()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", txID)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: txID,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()
	s := grpc.NewServer()
	service := &paymentService{}
	paymentV1.RegisterPaymentServiceServer(s, service)
	reflection.Register(s)
	go func() {
		log.Printf("🚀 gRPC server listening on port %d\n", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
