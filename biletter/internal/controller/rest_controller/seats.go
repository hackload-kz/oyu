package rest_controller

import (
	"biletter/internal/domain/entity"
	"encoding/json"
	"net/http"
)

func (h *docHandler) GetSeatsList(w http.ResponseWriter, r *http.Request) (err error) {

	var query entity.SeatsListQuery

	// получаем из query query параметр eds и валидируем
	err = query.Prepare(r.URL.Query())
	if err != nil {
		h.logger.Error(err)
		return err
	}

	var data []entity.SeatForList
	data, err = h.service.GetSeatsList(r.Context(), query)
	if err != nil {
		return err
	}

	var response []byte
	response, err = json.Marshal(data)
	if err != nil {
		h.logger.Error(err)
		return err
	}

	_, err = w.Write(response)
	if err != nil {
		h.logger.Error(err)
		return err
	}

	w.WriteHeader(http.StatusOK)

	h.logger.Info("Успешно получен список событий)")

	return nil
}
