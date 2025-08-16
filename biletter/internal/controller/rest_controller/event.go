package rest_controller

import (
	"biletter/internal/domain/entity"
	"encoding/json"
	"net/http"
)

// GetEventList
// @Summary Получение списка ивентов
// @Tags Event
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Description Получение списка ивентов
// @Accept json
// @Produce json
// @Success 201 {object} []entity.EventForList "..."
// @Failure 400,404 {object} error
// @Failure 500 {object} error
// @Failure default {object} error
// @Router /api/events [get]
func (h *docHandler) GetEventList(w http.ResponseWriter, r *http.Request) (err error) {

	var query entity.EventListQuery

	// получаем из query query параметр eds и валидируем
	err = query.Prepare(r.URL.Query())
	if err != nil {
		h.logger.Error(err)
		return err
	}

	var data []entity.EventForList
	data, err = h.service.GetEventList(r.Context(), query)
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
