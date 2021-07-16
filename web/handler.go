package web

import (
	"context"
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/connorjcantrell/toolint"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

func NewHandler(store toolint.Store, sessions *scs.SessionManager, csrfKey []byte) *Handler {
	h := &Handler{
		Mux:      chi.NewMux(),
		store:    store,
		sessions: sessions,
	}

	tools := ToolHandler{store: store, sessions: sessions}
	toolEntries := ToolEntryHandler{store: store, sessions: sessions}
	users := UserHandler{store: store, sessions: sessions}

	h.Use(middleware.Logger)
	h.Use(csrf.Protect(csrfKey, csrf.Secure(false)))
	h.Use(sessions.LoadAndSave)
	h.Use(h.withUser)

	h.Get("/", h.Home())
	h.Get("/tool/new", tools.Create())
	h.Post("/tool/new", tools.Store())
	h.Get("/tool-entry/new", toolEntries.Create())
	h.Post("/tool-entry/new", toolEntries.Store())
	h.Get("/register", users.Register())
	h.Post("/register", users.RegisterSubmit())
	h.Get("/login", users.Login())
	h.Post("/login", users.LoginSubmit())
	h.Get("/logout", users.Logout())

	return h
}

type Handler struct {
	*chi.Mux

	store    toolint.Store
	sessions *scs.SessionManager
}

func (h *Handler) Home() http.HandlerFunc {
	type data struct {
		SessionData

		ToolEntries []toolint.ToolEntry
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData := GetSessionData(h.sessions, r.Context())

		if !sessionData.LoggedIn {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		tt, err := h.store.ToolEntries()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{
			SessionData: sessionData,
			ToolEntries: tt,
		})
	}
}

func (h *Handler) withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := h.sessions.Get(r.Context(), "user_id").(uuid.UUID)

		user, err := h.store.User(id)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// TODO: Possible type collision
		// Define my own type (instead of user) to avoid collisions
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
