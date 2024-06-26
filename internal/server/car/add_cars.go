package car

import (
	"encoding/json"
	"net/http"

	"github.com/Uikola/carCatalog/internal/entity"
	"github.com/rs/zerolog/log"
)

// AddCars godoc
//
//	@Summary		Добавляет автомобили по их регистрационным номерам.
//	@Description	Получает на вход список регистрационных номеров, обогатив данные через стороннее api, добавляет их.
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Param			regNums	body		entity.AddCarsRequest	true	"Запрос добавления автомобиля"
//	@Success		201		{object}	entity.AddCarsResponse
//	@Failure		400		{object}	map[string]string
//	@Router			/cars [post]
func (h Handler) AddCars(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Adding cars")
	ctx := r.Context()

	var addCarsRequest entity.AddCarsRequest
	if err := json.NewDecoder(r.Body).Decode(&addCarsRequest); err != nil {
		log.Error().Err(err).Msg("bad json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad json"})
		return
	}

	addCarsRequest.Valid()

	cars := h.carUseCase.AddCars(ctx, addCarsRequest)

	log.Info().Msgf("Added %d cars", len(cars))

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(entity.AddCarsResponse{AddedCars: cars})
}
