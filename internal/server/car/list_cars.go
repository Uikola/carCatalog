package car

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Uikola/carCatalog/internal/entity"
	"github.com/Uikola/carCatalog/internal/errorz"
	"github.com/Uikola/carCatalog/pkg/filter"
	"github.com/rs/zerolog/log"
)

func (h Handler) ListCars(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var limit, offset int

	offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 5
	}

	filterOptions, err := FilterOptions(r)
	if err != nil {
		log.Error().Err(err).Msg("error while getting filter options")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "invalid options"})
		return
	}

	options := entity.ListCarsOptions{
		FilterOptions: filterOptions,
		Limit:         limit,
		Offset:        offset,
	}

	cars, err := h.carUseCase.ListCars(ctx, options)
	switch {
	case errors.Is(errorz.ErrNoCars, err):
		log.Debug().Msg("No cars found")
		_ = json.NewEncoder(w).Encode([]entity.Car{})
		w.WriteHeader(http.StatusOK)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("error while getting the list of cars")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "internal error"})
		return
	}

	log.Info().Msgf("Retrieved %d cars", len(cars))

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(cars)
}

func FilterOptions(r *http.Request) (Options, error) {
	log.Info().Msg("Listing cars")
	filterOptions := filter.NewOptions()

	log.Debug().Msg("Filter options:")

	regNum := r.URL.Query().Get("regNum")
	if regNum != "" {
		log.Debug().Msgf("regNum: %s", regNum)
		if err := filterOptions.AddField("reg_num", filter.OperatorLike, regNum, filter.DataTypeStr); err != nil {
			return nil, err
		}
	}

	mark := r.URL.Query().Get("mark")
	if mark != "" {
		log.Debug().Msgf("mark: %s", mark)
		if err := filterOptions.AddField("mark", filter.OperatorLike, mark, filter.DataTypeStr); err != nil {
			return nil, err
		}
	}

	model := r.URL.Query().Get("model")
	if model != "" {
		log.Debug().Msgf("model: %s", model)
		if err := filterOptions.AddField("model", filter.OperatorLike, model, filter.DataTypeStr); err != nil {
			return nil, err
		}
	}

	year := r.URL.Query().Get("year")
	if year != "" {
		operator := filter.OperatorEq
		val := year
		if strings.Index(year, ":") != -1 {
			split := strings.Split(year, ":")
			operator = split[0]
			val = split[1]
		}
		log.Debug().Msgf("year: %d", year)
		if err := filterOptions.AddField("year", operator, val, filter.DataTypeInt); err != nil {
			return nil, err
		}
	}

	ownerName := r.URL.Query().Get("owner_name")
	if ownerName != "" {
		log.Debug().Msgf("ownerName: %s", ownerName)
		if err := filterOptions.AddField("owner_name", filter.OperatorLike, ownerName, filter.DataTypeStr); err != nil {
			return nil, err
		}
	}

	ownerSurname := r.URL.Query().Get("owner_surname")
	if ownerSurname != "" {
		log.Debug().Msgf("ownerSurname: %s", ownerSurname)
		if err := filterOptions.AddField("owner_surname", filter.OperatorLike, ownerSurname, filter.DataTypeStr); err != nil {
			return nil, err
		}
	}

	ownerPatronymic := r.URL.Query().Get("owner_patronymic")
	if ownerPatronymic != "" {
		log.Debug().Msgf("ownerPatronymic: %s", ownerPatronymic)
		if err := filterOptions.AddField("owner_patronymic", filter.OperatorLike, ownerPatronymic, filter.DataTypeStr); err != nil {
			return nil, err
		}
	}

	return filterOptions, nil
}
