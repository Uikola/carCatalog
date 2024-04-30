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

// ListCars godoc
//
//	@Summary		Выводит список машин
//	@Description	Выводит список машин с фильтрацией по всем полям и пагинацией
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Param			limit				query		int		false	"Лимит"
//	@Param			offset				query		int		false	"Смещение"
//	@Param			regNum				query		string	false	"Фильтр по регистрационному номеру"
//	@Param			mark				query		string	false	"Фильтр по марке автомобиля"
//	@Param			model				query		string	false	"Фильтр по модели автомобиля"
//	@Param			mark				query		string	false	"Фильтр по марке автомобиля"
//	@Param			year				query		string	false	"Фильтр по году"	Format("year=2020 or year=gt:2020")
//	@Param			owner_name			query		string	false	"Фильтр по имени владельца"
//	@Param			owner_surname		query		string	false	"Фильтр по фамилии владельца"
//	@Param			owner_patronymic	query		string	false	"Фильтр по отчеству владельца"
//	@Success		200					{array}		entity.Car
//	@Failure		400					{object}	map[string]string
//	@Failure		500					{object}	map[string]string
//	@Router			/cars [get]
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
