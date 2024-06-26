package car

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Uikola/carCatalog/internal/entity"
	"github.com/Uikola/carCatalog/internal/errorz"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// UpdateCar godoc
//
//	@Summary		Обновляет поля автомобиля
//	@Description	Обновляет одно или несколько полей автомобиля
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Param			id				path	int						true	"Car ID"
//	@Param			carForUpdate	body	entity.UpdateCarRequest	true	"Запрос обновления автомобиля"
//	@Success		200				"OK"
//	@Failure		400				{object}	map[string]string
//	@Failure		404				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/cars/{id} [patch]
func (h Handler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	carID, err := strconv.Atoi(chi.URLParam(r, "car_id"))
	if err != nil {
		log.Error().Err(err).Msg("invalid car id")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad request"})
		return
	}
	log.Info().Msgf("Updating car with ID: %d", carID)

	var updateCarRequest entity.UpdateCarRequest
	if err = json.NewDecoder(r.Body).Decode(&updateCarRequest); err != nil {
		log.Error().Err(err).Msg("bad json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad json"})
		return
	}

	if err = updateCarRequest.Valid(); err != nil {
		log.Error().Err(err).Msg("bad request")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad request"})
		return
	}

	err = h.carUseCase.UpdateCar(ctx, uint(carID), updateCarRequest)
	switch {
	case errors.Is(errorz.ErrCarNotFound, err):
		log.Error().Err(err).Msg("car with this id not found")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "car with this id not found"})
		return
	case err != nil:
		log.Error().Err(err).Msg("error while updating the car")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "internal error"})
		return
	}

	log.Info().Msgf("Car with ID %d updated successfully", carID)

	w.WriteHeader(http.StatusOK)
}
