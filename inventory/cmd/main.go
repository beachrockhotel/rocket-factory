package main

import (
	"context"
	"fmt"
	inventory_v1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const grpcPort = 50051

type inventoryService struct {
	// Встроенная заготовка, которая реализует все методы с заглушками.
	inventory_v1.UnimplementedInventoryServiceServer
	// Мютекс на запись чтобы можно было взаимодействовать с методом
	mu sync.RWMutex
	// Поле в которое сохраняются данные, оно получает значение по ключу детали
	inv map[string]*inventory_v1.Part
}

func (s *inventoryService) GetPart(_ context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	// Лочим на чтение, чтобы защититься от одновременной записи
	s.mu.RLock()
	defer s.mu.RUnlock()
	// проходим по структурке инвентори сервиса и используем мапу и достаём из неё uuid с помощью req
	part, ok := s.inv[req.Uuid]
	// если не достали, то возвращаем код ошибки не найденной, с помощью библиотеки codes и строку ошибки и передаём элмент
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with uuid %s not found", req.Uuid)
	}
	// если всё нашлось, то возвращаем ответ
	return &inventory_v1.GetPartResponse{
		// Структура Part элемент part
		Part: part,
	}, nil
}

// проверяет есть ли хотя бы ОДИН общий тег
func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

// проверяет есть ли хотя бы ОДИН общий тег
func containsAny(list1, list2 []string) bool {
	for _, item := range list1 {
		if contains(list2, item) {
			return true
		}
	}
	return false
}

func containsCategory(list []inventory_v1.Category, item inventory_v1.Category) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

// Для вывода либо фильтрации деталей
func (s *inventoryService) ListParts(_ context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	// Получаем список элементов для фильтрации
	filter := req.GetFilter()
	// Если фильтр nil, значит клиент вообще не указал условия — возвращаем все детали.
	if filter == nil {
		// если фильтр вообще не передан — вернём всё сразу
		s.mu.RLock()
		defer s.mu.RUnlock()
		// Создаём слайс под все детали с capacity = количеству деталей, чтобы избежать лишних аллокаций при append.
		all := make([]*inventory_v1.Part, 0, len(s.inv))
		// Добавляем каждую деталь в слайс, чтобы вернуть их все
		for _, part := range s.inv {
			all = append(all, part)
		}
		return &inventory_v1.ListPartsResponse{Parts: all}, nil
	}
	// RLock, чтобы несколько раз не читало или не записывало одноверменно и чтобы не словили панику
	s.mu.RLock()
	defer s.mu.RUnlock()
	// не знаю
	var result []*inventory_v1.Part
	for _, part := range s.inv {
		// фильтр по UUID
		if len(filter.Uuids) > 0 && !contains(filter.Uuids, part.GetUuid()) {
			continue
		}
		// фильтр по имени
		if len(filter.Names) > 0 && !contains(filter.Names, part.GetName()) {
			continue
		}
		// фильтр по категории
		if len(filter.Categories) > 0 && !containsCategory(filter.Categories, part.GetCategory()) {
			continue
		}
		// фильтр по стране производителя
		if len(filter.ManufacturerCountries) > 0 {
			m := part.GetManufacturer()
			if m == nil || !contains(filter.ManufacturerCountries, m.GetCountry()) {
				continue
			}
		}
		// фильтр по тегам
		if len(filter.Tags) > 0 && !containsAny(filter.Tags, part.GetTags()) {
			continue
		}
		result = append(result, part)
	}

	return &inventory_v1.ListPartsResponse{Parts: result}, nil
}

func main() {
	// устанавливает tcp соеденение на порту 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()
	// Создаёт grpc сервер
	s := grpc.NewServer()
	// Регистрируем наш сервис
	service := &inventoryService{
		inv: make(map[string]*inventory_v1.Part),
	}
	service.initParts()
	// Применяет методы из сгенерированного кода к сервису
	inventory_v1.RegisterInventoryServiceServer(s, service)
	// Включает рефлексию (
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
