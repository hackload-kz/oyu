package postgresql

import (
	"context"
)

func (s *storage) BookingCreate(ctx context.Context, eventID, userID int64) (bookingID int64, err error) {
	createQuery := `
		INSERT INTO
		    bookings (user_id, event_id) 
		VALUES ($1, $2)
		RETURNING id;
	`
	err = s.client.QueryRow(ctx,
		createQuery,
		userID,
		eventID,
	).Scan(&bookingID)

	if err != nil {
		s.logger.Errorf("не удалось добавить бронирование: ERROR=%v", err)
		return 0, err
	}

	return bookingID, nil
}
