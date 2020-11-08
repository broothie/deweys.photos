package app

import (
	"net/http"
	"sort"
)

func (a *App) Index() http.HandlerFunc {
	page := Page("index.html")

	type Data struct {
		SessionData
		Entries []Entry
	}

	return func(w http.ResponseWriter, r *http.Request) {
		entries, err := a.getAllEntries(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sort.Slice(entries, func(i, j int) bool { return entries[i].CreatedAt.After(entries[j].CreatedAt) })
		if err := page.Execute(w, Data{Entries: entries, SessionData: a.sessionData(w, r)}); err != nil {
			a.Logger.Println(err)
		}
	}
}
