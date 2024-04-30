package server

import (
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"

	"github.com/Uikola/carCatalog/internal/server/car"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/Uikola/carCatalog/docs"
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
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"), //The url pointing to API definition
	))

	router.Route("/api", func(r chi.Router) {

		r.Get("/cars", carHandler.ListCars)
		r.Delete("/cars/{car_id}", carHandler.DeleteCar)
		r.Patch("/cars/{car_id}", carHandler.UpdateCar)
		r.Post("/cars", carHandler.AddCars)
	})
}
