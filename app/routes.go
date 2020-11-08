package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a *App) Routes(router *mux.Router) {
	// File server
	router.
		Methods(http.MethodGet).
		PathPrefix("/public").
		Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Index
	router.
		Methods(http.MethodGet).
		Path("/").
		Handler(a.Index())

	// Login
	router.
		Methods(http.MethodGet).
		Path("/login").
		Handler(a.NewLogin())

	router.
		Methods(http.MethodPost).
		Path("/login").
		Handler(a.CreateLogin())

	router.
		Methods(http.MethodDelete).
		Path("/login").
		Handler(a.DeleteLogin())

	// Photos
	router.
		Methods(http.MethodGet).
		Path("/photos/manage").
		Handler(a.RequireLoggedIn(a.ManagePhotos()))

	router.
		Methods(http.MethodPost).
		Path("/photos").
		Handler(a.RequireLoggedIn(a.CreatePhoto()))

	router.
		Methods(http.MethodGet).
		Path("/photos/{entry_id}").
		Handler(a.RequireLoggedIn(a.EditPhoto()))

	router.
		Methods(http.MethodPost).
		Path("/photos/{entry_id}").
		Handler(a.RequireLoggedIn(a.UpdatePhoto()))
}
