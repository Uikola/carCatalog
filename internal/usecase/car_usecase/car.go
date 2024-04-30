package car_usecase

import (
	"context"
	"github.com/Uikola/carCatalog/internal/entity"
	"github.com/Uikola/carCatalog/internal/errorz"
	"github.com/Uikola/carCatalog/pkg/enrichment_api"
	"github.com/rs/zerolog/log"
)

func (uc UseCaseImpl) ListCars(ctx context.Context, options entity.ListCarsOptions) ([]entity.Car, error) {
	return uc.carRepository.ListCars(ctx, options)
}

func (uc UseCaseImpl) DeleteCar(ctx context.Context, carID uint) error {
	return uc.carRepository.DeleteCar(ctx, carID)
}

func (uc UseCaseImpl) UpdateCar(ctx context.Context, carID uint, request entity.UpdateCarRequest) error {
	exists, err := uc.carRepository.ExistsByID(ctx, carID)
	if err != nil {
		return err
	}

	if !exists {
		return errorz.ErrCarNotFound
	}

	carMap := map[string]interface{}{
		"reg_num":          request.RegNum,
		"mark":             request.Mark,
		"model":            request.Model,
		"year":             request.Year,
		"owner_name":       request.Owner.Name,
		"owner_surname":    request.Owner.Surname,
		"owner_patronymic": request.Owner.Patronymic,
	}

	return uc.carRepository.UpdateCar(ctx, carID, carMap)
}

func (uc UseCaseImpl) AddCars(ctx context.Context, request entity.AddCarsRequest) []entity.Car {
	client := enrichment_api.NewClient()

	var carsToAdd []entity.Car
	for _, regNum := range request.RegNums {
		exists, err := uc.carRepository.ExistsByRegNum(ctx, regNum)
		if err != nil {
			log.Error().Err(err).Msg("failed to check if exists car by registration number")
			continue
		}

		if exists {
			log.Error().Err(errorz.ErrCarAlreadyExists).Msg("car with this registration number already exists")
			continue
		}

		car, err := client.GetCarByRegNum(regNum)
		if err != nil {
			log.Error().Err(err).Msg("failed to get car by reg num")
			continue
		}
		carsToAdd = append(carsToAdd, car)
	}

	var cars []entity.Car
	for _, carToAdd := range carsToAdd {
		car, err := uc.carRepository.AddCar(ctx, carToAdd)
		if err != nil {
			log.Error().Err(err).Msg("failed to add car")
			continue
		}
		cars = append(cars, car)
	}

	if cars == nil {
		return []entity.Car{}
	}

	return cars
}
