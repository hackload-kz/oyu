package services

import (
	"biletter/internal/domain/entity"
	"context"
)

func (s *service) GetEventList(ctx context.Context, query entity.EventListQuery) (
	response []entity.EventForList, err error,
) {
	filters := query.GenerateFilters()

	response, err = s.postgresStorage.GetEventList(ctx, filters)
	if err != nil {
		return nil, err
	}

	return response, nil
}
