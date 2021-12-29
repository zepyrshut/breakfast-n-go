package main

import (
	"net/http"

	"github.com/zepyrshut/breakfast-n-go/internal/config"
	"github.com/zepyrshut/breakfast-n-go/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/availability", handlers.Repo.Availability)
	mux.Post("/availability", handlers.Repo.DoPostAvailability)
	mux.Post("/availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/rooms", handlers.Repo.Rooms)
	mux.Get("/terrace", handlers.Repo.Terrace)
	mux.Get("/spa", handlers.Repo.Spa)
	mux.Get("/reservation", handlers.Repo.Reservation)
	mux.Get("/contact", handlers.Repo.Contact)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
