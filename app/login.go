package app

import (
	"net/http"
	"os"
	"time"
)

const loggedInCookieName = "logged_in"

func (a *App) NewLogin() http.HandlerFunc {
	page := Page("login.html")

	return func(w http.ResponseWriter, r *http.Request) {
		if a.isLoggedIn(r) {
			http.Redirect(w, r, "/photos/manage", http.StatusPermanentRedirect)
			return
		}

		if err := page.Execute(w, a.sessionData(w, r)); err != nil {
			a.Logger.Println(err)
		}
	}
}

func (a *App) CreateLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("password") != os.Getenv("PASSWORD") {
			a.flash(w, "incorrect password")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		a.setCookie(w, loggedInCookieName, "true", 24*time.Hour)
		a.callback(w, r)
	}
}

func (a *App) DeleteLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.takeCookie(w, r, loggedInCookieName)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (a *App) RequireLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.isLoggedIn(r) {
			a.callbackRedirect(w, r, "/login")
			return
		}

		next.ServeHTTP(w, r)
	})
}
