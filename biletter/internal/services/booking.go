package services

import (
	"biletter/internal/domain/entity"
	"context"
	"errors"
)

func (s *service) BookingCreate(ctx context.Context, data entity.BookingCreate) (bookingID int64, err error) {
	user, ok := entity.UserFromContext(ctx)
	if !ok {
		err = errors.New("не удалось получить пользователя из контекста")
		s.logger.Error(err)
		return 0, err
	}

	// создание записи нового города
	bookingID, err = s.postgresStorage.BookingCreate(ctx, data.EventID, user.ID)
	if err != nil {
		return 0, err
	}

	return bookingID, nil
}
