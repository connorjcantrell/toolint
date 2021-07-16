package web

import (
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/connorjcantrell/toolint"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

type ToolEntryHandler struct {
	store    toolint.Store
	sessions *scs.SessionManager
}

func (h *ToolEntryHandler) Create() http.HandlerFunc {
	type data struct {
		SessionData

		CSRF      template.HTML
		ToolEntry toolint.ToolEntry
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/tool_entry_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := h.store.ToolEntry(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{
			SessionData: GetSessionData(h.sessions, r.Context()),
			CSRF:        csrf.TemplateField(r),
			ToolEntry:   t,
		})
	}
}

func (h *ToolEntryHandler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		toolID, err := uuid.Parse(r.FormValue("id"))
		if err != nil {
			h.sessions.Put(r.Context(), "flash", "Invalid Tool selected")
			http.Redirect(w, r, "/home/", http.StatusFound)
			return
		}
		form := CreateToolEntryForm{
			ToolID:    toolID,
			Condition: r.FormValue("condition"),
		}
		if !form.Validate() {
			h.sessions.Put(r.Context(), "form", form)
			w.Header().Set("Cache-Control", "private")
			http.Redirect(w, r, r.Referer(), http.StatusFound)
			return
		}

		data := GetSessionData(h.sessions, r.Context())

		t := &toolint.ToolEntry{
			ID:        uuid.New(),
			ToolID:    form.ToolID,
			UserID:    data.User.ID,
			Condition: form.Condition,
		}
		if err := h.store.CreateToolEntry(t); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.sessions.Put(r.Context(), "flash", "Tool Entry has been created.")

		http.Redirect(w, r, "/home/", http.StatusFound)
	}
}

func (h *ToolEntryHandler) Show() http.HandlerFunc {
	type data struct {
		SessionData
		CSRF      template.HTML
		ToolEntry toolint.ToolEntry
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/tool_entry.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		toolEntryIDStr := chi.URLParam(r, "tool_entry_id")

		toolEntryID, err := uuid.Parse(toolEntryIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := h.store.ToolEntry(toolEntryID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{
			SessionData: GetSessionData(h.sessions, r.Context()),
			CSRF:        csrf.TemplateField(r),
			ToolEntry:   t,
		})
	}
}

func (h *ToolEntryHandler) List() http.HandlerFunc {
	type data struct {
		SessionData
		ToolEntries []toolint.ToolEntry
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/tool_entries.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tt, err := h.store.ToolEntries()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{
			SessionData: GetSessionData(h.sessions, r.Context()),
			ToolEntries: tt,
		})
	}
}
