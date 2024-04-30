package car_usecase

import (
	"context"
	"github.com/Uikola/carCatalog/internal/entity"
	"github.com/Uikola/carCatalog/pkg/filter"
)

type Options interface {
	IsToApply() bool
	AddField(name, operator, value, dType string) error
	Fields() []filter.Field
}

type carRepository interface {
	ListCars(ctx context.Context, options entity.ListCarsOptions) ([]entity.Car, error)
	DeleteCar(ctx context.Context, carID uint) error
	UpdateCar(ctx context.Context, carID uint, carMap map[string]interface{}) error
	ExistsByID(ctx context.Context, carID uint) (bool, error)
	ExistsByRegNum(ctx context.Context, regNum string) (bool, error)
	AddCar(ctx context.Context, car entity.Car) (entity.Car, error)
}

type UseCaseImpl struct {
	carRepository carRepository
}

func NewUseCaseImpl(carRepository carRepository) *UseCaseImpl {
	return &UseCaseImpl{
		carRepository: carRepository,
	}
}
