package services

import (
	"biletter/internal/domain/entity"
	"context"
	"fmt"
	"github.com/go-redis/cache/v9"
	"time"
)

func (s *service) GetSeatsList(ctx context.Context, query entity.SeatsListQuery) (
	response []entity.SeatForList, err error,
) {
	cacheKey := query.GetCacheKey()

	// 1) Попытка из Redis
	err = s.redisCache.Get(ctx, cacheKey, &response)
	if err == nil || len(response) != 0 {
		return response, nil
	}

	response, err = s.postgresStorage.GetSeatsList(ctx, query)
	if err != nil {
		return nil, err
	}

	err = s.redisCache.Set(&cache.Item{
		Ctx:   context.Background(),
		Key:   cacheKey,
		Value: response,
		TTL:   20 * time.Second,
	})
	if err != nil {
		err = fmt.Errorf("не удалось положить новые значение списка мест в кеш под ключом=%s : ERROR=%v", cacheKey, err)
		s.logger.Error(err)
	}

	return response, nil
}
