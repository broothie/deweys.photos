package app

import (
	"net/http"
	"time"
)

const (
	flashCookieName    = "flash"
	redirectCookieName = "callback"
)

type SessionData struct {
	LoggedIn bool
	Flash    string
}

func (a *App) flash(w http.ResponseWriter, message string) {
	a.setCookie(w, flashCookieName, message, 5*time.Second)
}

func (a *App) isLoggedIn(r *http.Request) bool {
	return a.getCookie(r, loggedInCookieName) != ""
}

func (a *App) sessionData(w http.ResponseWriter, r *http.Request) SessionData {
	return SessionData{
		LoggedIn: a.isLoggedIn(r),
		Flash:    a.takeCookie(w, r, flashCookieName),
	}
}

func (a *App) callbackRedirect(w http.ResponseWriter, r *http.Request, url string) {
	a.setCookie(w, redirectCookieName, r.URL.Path, 5*time.Minute)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (a *App) callback(w http.ResponseWriter, r *http.Request) {
	redirectPath := a.takeCookie(w, r, redirectCookieName)
	if redirectPath == "" {
		redirectPath = "/"
	}

	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
}

func (a *App) getCookie(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		if err != http.ErrNoCookie {
			a.Logger.Println("failed to retrieve cookie:", name)
		}

		return ""
	}

	return cookie.Value
}

func (a *App) takeCookie(w http.ResponseWriter, r *http.Request, name string) string {
	value := a.getCookie(r, name)
	http.SetCookie(w, &http.Cookie{Name: name, MaxAge: -1})
	return value
}

func (a *App) setCookie(w http.ResponseWriter, name, value string, expires time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(expires),
		Secure:   a.Environment == "production",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
