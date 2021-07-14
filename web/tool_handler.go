package web

import (
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	toolint "github.com/connorjcantrell/toolint"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

type ToolHandler struct {
	store    toolint.Store
	sessions *scs.SessionManager
}

func (h *ToolHandler) Create() http.HandlerFunc {
	type data struct {
		SessionData

		CSRF template.HTML
		Tool toolint.Tool
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/tool_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := h.store.Tool(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{
			SessionData: GetSessionData(h.sessions, r.Context()),
			CSRF:        csrf.TemplateField(r),
			Tool:        t,
		})
	}
}

func (h *ToolHandler) Show() http.HandlerFunc {
	type data struct {
		SessionData
		CSRF template.HTML
		Tool toolint.Tool
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/tool.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		toolIDStr := chi.URLParam(r, "toolID")

		toolID, err := uuid.Parse(toolIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := h.store.Tool(toolID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{
			SessionData: GetSessionData(h.sessions, r.Context()),
			CSRF:        csrf.TemplateField(r),
			Tool:        t,
		})
	}
}
