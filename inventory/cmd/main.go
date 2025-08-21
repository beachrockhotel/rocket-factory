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

var ErrPartNotFound = errors.New("part not found")

type inventoryService struct {
	inventory_v1.UnimplementedInventoryServiceServer
	mu  sync.RWMutex
	inv map[string]*inventory_v1.Part
}

func (s *inventoryService) GetPart(_ context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	part, ok := s.inv[req.Uuid]
	if !ok {
		return nil, ErrPartNotFound
	}
	return &inventory_v1.GetPartResponse{
		Part: part,
	}, nil
}

func isEmptyFilter(f *inventory_v1.PartsFilter) bool {
	if f == nil {
		return true
	}
	return len(f.Uuids) == 0 &&
		len(f.Names) == 0 &&
		len(f.Categories) == 0 &&
		len(f.ManufacturerCountries) == 0 &&
		len(f.Tags) == 0
}

func hasAnyTag(need, have []string) bool {
	if len(need) == 0 || len(have) == 0 {
		return false
	}
	set := make(map[string]struct{}, len(have))
	for _, t := range have {
		set[t] = struct{}{}
	}
	for _, t := range need {
		if _, ok := set[t]; ok {
			return true
		}
	}
	return false
}

func (s *inventoryService) ListParts(_ context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	f := req.GetFilter()

	s.mu.RLock()
	defer s.mu.RUnlock()

	if isEmptyFilter(f) {
		all := make([]*inventory_v1.Part, 0, len(s.inv))
		for _, part := range s.inv {
			all = append(all, part)
		}
		return &inventory_v1.ListPartsResponse{Parts: all}, nil
	}

	result := make([]*inventory_v1.Part, 0, 16)
	for _, part := range s.inv {
		uuidOK := len(f.Uuids) == 0 || slices.Contains(f.Uuids, part.GetUuid())
		nameOK := len(f.Names) == 0 || slices.Contains(f.Names, part.GetName())
		catOK := len(f.Categories) == 0 || slices.Contains(f.Categories, part.GetCategory())

		countryOK := true
		if len(f.ManufacturerCountries) > 0 {
			m := part.GetManufacturer()
			countryOK = m != nil && slices.Contains(f.ManufacturerCountries, m.GetCountry())
		}

		tagsOK := len(f.Tags) == 0 || hasAnyTag(f.Tags, part.GetTags())

		if uuidOK && nameOK && catOK && countryOK && tagsOK {
			result = append(result, part)
		}
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
