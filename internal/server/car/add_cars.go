package car

import (
	"encoding/json"
	"github.com/Uikola/carCatalog/internal/entity"
	"github.com/rs/zerolog/log"
	"net/http"
)

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
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"added_cars": cars})
}
