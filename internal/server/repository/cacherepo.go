package repository

import (
	"errors"
	"wb_level_0/internal/cache"
	"wb_level_0/internal/order"
)

type CacheRepository struct {
	cache *cache.Cache
}

func NewCacheRepository() *CacheRepository {

	return &CacheRepository{
		cache: cache.NewCache(),
	}
}

func (repo *CacheRepository) GetById(uid string) (*order.Order, error) {
	result, ok := repo.cache.Get(uid)
	if !ok {
		return nil, errors.New("no order with such UID in cache")
	}
	return result.(*order.Order), nil
}

func (repo *CacheRepository) Set(uid string, order *order.Order) error {
	repo.cache.Set(order.OrderUID, order)
	return nil
}

func (repo *CacheRepository) IsEmpty() bool {
	return repo.cache.Len() == 0
}