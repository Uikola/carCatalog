package car

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// DeleteCar godoc
//
//	@Summary		Удаляет автомобиль
//	@Description	Удаляет автомобиль по его ID
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Car ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/cars/{id} [delete]
func (h Handler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	carID, err := strconv.Atoi(chi.URLParam(r, "car_id"))
	if err != nil {
		log.Error().Err(err).Msg("invalid car id")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad request"})
		return
	}
	log.Info().Msgf("Deleting car with ID: %d", carID)

	err = h.carUseCase.DeleteCar(ctx, uint(carID))
	if err != nil {
		log.Error().Err(err).Msg("error while deleting the car")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "internal error"})
		return
	}

	log.Info().Msgf("Car with ID %d deleted successfully", carID)

	w.WriteHeader(http.StatusNoContent)
}
