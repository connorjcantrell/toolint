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

func (h *ToolHandler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := CreateToolForm{
			Name:     r.FormValue("name"),
			Model:    r.FormValue("model"),
			Category: r.FormValue("category"),
		}
		if !form.Validate() {
			h.sessions.Put(r.Context(), "form", form)
			http.Redirect(w, r, r.Referer(), http.StatusFound)
			return
		}

		t := &toolint.Tool{
			ID:       uuid.New(),
			Name:     form.Name,
			Model:    form.Model,
			Category: form.Category,
		}
		if err := h.store.CreateTool(t); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.sessions.Put(r.Context(), "flash", "Tool has been created.")

		http.Redirect(w, r, "/tools/", http.StatusFound)
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
		toolIDStr := chi.URLParam(r, "ID")

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

func (h *ToolHandler) List() http.HandlerFunc {
	type data struct {
		SessionData
		Tools []toolint.Tool
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/tools.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tt, err := h.store.Tools()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{
			SessionData: GetSessionData(h.sessions, r.Context()),
			Tools:       tt,
		})
	}
}
