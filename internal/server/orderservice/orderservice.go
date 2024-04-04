package orderservice

import (
	"fmt"
	"wb_level_0/internal/order"
	"wb_level_0/internal/server/repository"
)

type OrderService struct {
	cacheRepository *repository.CacheRepository
	sqlRepository   *repository.SqlRepository
}

func NewOrderService(cacheRepository *repository.CacheRepository, sqlRepository *repository.SqlRepository) *OrderService {
	return &OrderService{
		cacheRepository: cacheRepository,
		sqlRepository:   sqlRepository,
	}
}

func (s *OrderService) GetByUid(orderUid string) (*order.Order, error) {
	if !s.cacheRepository.IsEmpty() {
		return s.cacheRepository.GetById(orderUid)
	}
	return nil, fmt.Errorf("empty cache")
}

func (s *OrderService) CacheRecovery() error {
	orders, err := s.sqlRepository.GetAllOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		s.cacheRepository.Set(order.OrderUID, order)
	}
	return nil
}

func (s *OrderService) SaveOrder(order *order.Order) error {
	err := s.sqlRepository.SaveOrder(order)
	if err != nil {
		return err
	}
	err = s.cacheRepository.Set(order.OrderUID, order)
	return err
}
