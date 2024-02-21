package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"dev.acmcsuf.com/fullyhacks-qrms/frontend"
	"dev.acmcsuf.com/fullyhacks-qrms/sqldb"
	"github.com/go-chi/chi/v5"
	"libdb.so/tmplutil"
)

// Handler is the type for handling API requests.
type Handler struct {
	mux    *chi.Mux
	db     *sqldb.Database
	tmpl   *tmplutil.Templater
	logger *slog.Logger
}

// NewHandler creates a new Handler.
func NewHandler(db *sqldb.Database, logger *slog.Logger) *Handler {
	h := &Handler{
		mux:    chi.NewMux(),
		db:     db,
		tmpl:   frontend.NewTemplater(),
		logger: logger,
	}

	h.tmpl.OnRenderFail = func(sub *tmplutil.Subtemplate, w io.Writer, err error) {
		h.renderError(w, err)
	}

	m := h.mux
	m.NotFound(func(w http.ResponseWriter, r *http.Request) {
		h.renderErrorWithCode(w, 404, fmt.Errorf("page not found"))
	})

	m.Get("/", h.getIndex)

	m.Route("/events", func(m chi.Router) {
		m.Get("/new", h.getNewEvent)
		m.Post("/new", h.postNewEvent)

		m.Route("/{id}", func(m chi.Router) {
			m.Get("/", h.getEvent)
			m.Get("/attendees", h.getEventAttendees)

			m.Get("/scan", h.getScanEvent)
			m.Post("/scan", h.postScanEvent)

			m.Get("/merge", h.getMergeEvent)
			m.Post("/merge", h.postMergeEvent)
		})
	})

	m.Mount("/", frontend.StaticHandler())

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) renderErrorWithCode(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	h.renderError(w, err)
}

func (h *Handler) renderError(w io.Writer, err error) {
	if renameErr := h.tmpl.Execute(w, "error", err); renameErr != nil {
		h.logger.Error(
			"Failed to render error page",
			"render_error", renameErr,
			"actual_error", err)
	}
}

type indexData struct {
	Events []sqldb.ListEventsRow
}

func (h *Handler) getIndex(w http.ResponseWriter, r *http.Request) {
	events, err := h.db.ListEvents(r.Context())
	if err != nil {
		h.renderErrorWithCode(w, 500, fmt.Errorf("cannot get events: %w", err))
		return
	}

	h.tmpl.Execute(w, "index", indexData{
		Events: events,
	})
}

func (h *Handler) getNewEvent(w http.ResponseWriter, r *http.Request) {
	h.tmpl.Execute(w, "event_new", nil)
}

func (h *Handler) postNewEvent(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	if description == "" {
		h.renderErrorWithCode(w, 400, fmt.Errorf("description cannot be empty"))
		return
	}

	uuid := sqldb.GenerateUUID()

	if err := h.db.CreateEvent(r.Context(), sqldb.CreateEventParams{
		UUID:        uuid,
		Description: description,
	}); err != nil {
		h.renderErrorWithCode(w, 500, fmt.Errorf("cannot create event: %w", err))
		return
	}

	http.Redirect(w, r, "/events/"+uuid, http.StatusSeeOther)
}

func (h *Handler) getEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	slog.Debug(
		"Getting event",
		"id", id)

	event, err := h.db.GetEvent(r.Context(), id)
	if err != nil {
		h.renderErrorWithCode(w, 500, fmt.Errorf("cannot get event: %w", err))
		return
	}

	h.tmpl.Execute(w, "event", event)
}

func (h *Handler) getScanEvent(w http.ResponseWriter, r *http.Request) {
	h.tmpl.Execute(w, "scan", nil)
}

func (h *Handler) postScanEvent(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) getEventAttendees(w http.ResponseWriter, r *http.Request) {
	h.tmpl.Execute(w, "event_attendees", nil)
}

func (h *Handler) getMergeEvent(w http.ResponseWriter, r *http.Request) {
	h.tmpl.Execute(w, "event_merge", nil)
}

func (h *Handler) postMergeEvent(w http.ResponseWriter, r *http.Request) {}
