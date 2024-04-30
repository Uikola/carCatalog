package server

import (
	"github.com/Uikola/carCatalog/internal/server/car"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func NewServer(
	carHandler *car.Handler,
) http.Handler {
	router := chi.NewRouter()

	addRoutes(router, carHandler)

	var handler http.Handler = router

	return handler
}

func addRoutes(
	router *chi.Mux,
	carHandler *car.Handler,
) {
	router.Use(middleware.Logger)

	router.Route("/api", func(r chi.Router) {
		r.Get("/cars", carHandler.ListCars)
		r.Delete("/cars/{car_id}", carHandler.DeleteCar)
		r.Patch("/cars/{car_id}", carHandler.UpdateCar)
		r.Post("/cars", carHandler.AddCars)
	})
}
