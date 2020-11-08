package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (a *App) ManagePhotos() http.HandlerFunc {
	page := Page("manage_photos.html")

	type Data struct {
		SessionData
		Entry   Entry
		Entries []Entry
		Paths   map[string]string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		entries, err := a.getAllEntries(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := Data{
			SessionData: a.sessionData(w, r),
			Entries:     entries,
			Paths:       map[string]string{"action": "/photos", "cancel": "/"},
		}

		if err := page.Execute(w, data); err != nil {
			a.Logger.Println(err)
		}
	}
}

func (a *App) CreatePhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uploadedFile, handler, err := r.FormFile("photo")
		if err != nil {
			a.Logger.Println(err)
			a.flash(w, "missing photo")
			http.Redirect(w, r, "/photos/manage", http.StatusSeeOther)
			return
		}
		defer uploadedFile.Close()

		cloudFile := a.Bucket.Object(handler.Filename).NewWriter(r.Context())
		if _, err := io.Copy(cloudFile, uploadedFile); err != nil {
			a.Logger.Println(err)
			a.flash(w, "failed to copy photo contents")
			http.Redirect(w, r, "/photos/manage", http.StatusSeeOther)
			return
		}

		if err := cloudFile.Close(); err != nil {
			a.Logger.Println(err)
			a.flash(w, "failed to write file to cloud")
			http.Redirect(w, r, "/photos/manage", http.StatusSeeOther)
			return
		}

		attrs := cloudFile.Attrs()
		_, _, err = a.Entries.Add(r.Context(), Entry{
			Title:     r.FormValue("title"),
			Blurb:     r.FormValue("blurb"),
			Filename:  attrs.Name,
			URL:       attrs.MediaLink,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			a.Logger.Println(err)
			a.flash(w, "failed to save photo details")
			http.Redirect(w, r, "/photos/manage", http.StatusSeeOther)
			return
		}

		a.flash(w, fmt.Sprintf("%s uploaded", handler.Filename))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (a *App) EditPhoto() http.HandlerFunc {
	page := Page("edit_photo.html")

	type Data struct {
		SessionData
		Entry Entry
	}

	return func(w http.ResponseWriter, r *http.Request) {
		entryID := mux.Vars(r)["entry_id"]
		doc, err := a.Entries.Doc(entryID).Get(r.Context())
		if err != nil {
			a.Logger.Println(err)
			a.flash(w, "failed to find photo")
			http.Redirect(w, r, "/photos/manage", http.StatusSeeOther)
			return
		}

		var entry Entry
		if err := doc.DataTo(&entry); err != nil {
			a.Logger.Println(err)
			a.flash(w, "failed to read photo data")
			http.Redirect(w, r, "/photos/manage", http.StatusSeeOther)
			return
		}

		entry.ID = doc.Ref.ID
		if err := page.Execute(w, Data{Entry: entry, SessionData: a.sessionData(w, r)}); err != nil {
			a.Logger.Println(err)
		}
	}
}

func (a *App) UpdatePhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entryID := mux.Vars(r)["entry_id"]
		_, err := a.Entries.Doc(entryID).Update(r.Context(), []firestore.Update{
			{Path: "title", Value: r.FormValue("title")},
			{Path: "blurb", Value: r.FormValue("blurb")},
		})
		if err != nil {
			a.Logger.Println(err)
			a.flash(w, "failed to update photo data")
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/#photo-%s", entryID), http.StatusSeeOther)
	}
}

func (a *App) getAllEntries(ctx context.Context) ([]Entry, error) {
	docs, err := a.Entries.Documents(ctx).GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get entries")
	}

	entries := make([]Entry, 0, len(docs))
	for _, doc := range docs {
		var entry Entry
		if err := doc.DataTo(&entry); err != nil {
			a.Logger.Println(err)
			continue
		}

		entry.ID = doc.Ref.ID
		entries = append(entries, entry)
	}

	return entries, nil
}
