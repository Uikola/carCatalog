package car

import (
	"context"

	"github.com/Uikola/carCatalog/internal/entity"
	"github.com/Uikola/carCatalog/pkg/filter"
)

type carUseCase interface {
	ListCars(ctx context.Context, options entity.ListCarsOptions) ([]entity.Car, error)
	DeleteCar(ctx context.Context, carID uint) error
	UpdateCar(ctx context.Context, carID uint, request entity.UpdateCarRequest) error
	AddCars(ctx context.Context, request entity.AddCarsRequest) []entity.Car
}

type Options interface {
	IsToApply() bool
	AddField(name, operator, value, dType string) error
	Fields() []filter.Field
}

type Handler struct {
	carUseCase carUseCase
}

func NewHandler(carUseCase carUseCase) *Handler {
	return &Handler{
		carUseCase: carUseCase,
	}
}
