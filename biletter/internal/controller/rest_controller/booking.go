package rest_controller

import (
	"biletter/internal/domain/entity"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func (h *docHandler) BookingCreate(w http.ResponseWriter, r *http.Request) (err error) {
	// 1) Ограничим тело (например, 8KB) и запретим лишние поля
	r.Body = http.MaxBytesReader(w, r.Body, 8<<10)
	defer r.Body.Close()

	var booking entity.BookingCreate

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err = dec.Decode(&booking); err != nil {
		h.logger.Error(err)
		return err
	}
	if err = booking.Validate(); err != nil {
		h.logger.Error(err)
		return err
	}

	// тайм-аут на хендлер (чтобы не висеть под нагрузкой)
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	bookingID, err := h.service.BookingCreate(ctx, booking)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(
		struct {
			ID int64 `json:"id"`
		}{ID: bookingID})

	h.logger.Infof("Успешно добавлен бронирование: ID=%d", bookingID)

	return nil
}
