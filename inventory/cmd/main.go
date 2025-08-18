package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	inventory_v1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

type inventoryService struct {
	inventory_v1.UnimplementedInventoryServiceServer
	mu  sync.RWMutex
	inv map[string]*inventory_v1.Part
}

func (s *inventoryService) GetPart(_ context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var ErrPartNotFound = errors.New("part not found")
	part, ok := s.inv[req.Uuid]
	if !ok {
		return nil, ErrPartNotFound
	}
	return &inventory_v1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *inventoryService) ListParts(_ context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	filter := req.GetFilter()

	s.mu.RLock()
	defer s.mu.RUnlock()

	if filter == nil ||
		(len(filter.Uuids) == 0 &&
			len(filter.Names) == 0 &&
			len(filter.Categories) == 0 &&
			len(filter.ManufacturerCountries) == 0 &&
			len(filter.Tags) == 0) {
		all := make([]*inventory_v1.Part, 0, len(s.inv))
		for _, part := range s.inv {
			all = append(all, part)
		}
		return &inventory_v1.ListPartsResponse{Parts: all}, nil
	}

	result := make([]*inventory_v1.Part, 0, 16)
	for _, part := range s.inv {
		if len(filter.Uuids) > 0 && !slices.Contains(filter.Uuids, part.GetUuid()) {
			continue
		}
		if len(filter.Names) > 0 && !slices.Contains(filter.Names, part.GetName()) {
			continue
		}
		if len(filter.Categories) > 0 && !slices.Contains(filter.Categories, part.GetCategory()) {
			continue
		}
		if len(filter.ManufacturerCountries) > 0 {
			m := part.GetManufacturer()
			if m == nil || !slices.Contains(filter.ManufacturerCountries, m.GetCountry()) {
				continue
			}
		}
		if len(filter.Tags) > 0 {
			ptags := part.GetTags()
			hasAny := false
			for _, t := range filter.Tags {
				if slices.Contains(ptags, t) {
					hasAny = true
					break
				}
			}
			if !hasAny {
				continue
			}
		}

		result = append(result, part)
	}

	return &inventory_v1.ListPartsResponse{Parts: result}, nil
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
	service := &inventoryService{
		inv: make(map[string]*inventory_v1.Part),
	}
	service.initParts()
	inventory_v1.RegisterInventoryServiceServer(s, service)
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on port %d\n", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
