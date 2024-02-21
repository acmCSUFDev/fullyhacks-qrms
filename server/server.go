package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"dev.acmcsuf.com/fullyhacks-qrms/frontend"
	"dev.acmcsuf.com/fullyhacks-qrms/sqldb"
	"libdb.so/tmplutil"
)

// Handler is the type for handling API requests.
type Handler struct {
	mux    *http.ServeMux
	db     *sqldb.Database
	tmpl   *tmplutil.Templater
	logger *slog.Logger
}

// NewHandler creates a new Handler.
func NewHandler(db *sqldb.Database, logger *slog.Logger) *Handler {
	h := &Handler{
		mux:    http.NewServeMux(),
		db:     db,
		tmpl:   frontend.NewTemplater(),
		logger: logger,
	}

	h.tmpl.OnRenderFail = func(sub *tmplutil.Subtemplate, w io.Writer, err error) {
		h.renderError(w, err)
	}

	m := h.mux

	m.HandleFunc("GET /{$}", h.getIndex)

	m.HandleFunc("GET /events/new/{$}", h.getNewEvent)
	m.HandleFunc("POST /events/new/{$}", h.postNewEvent)

	m.HandleFunc("GET /events/{id}/{$}", h.getEvent)
	m.HandleFunc("POST /events/{id}/scan/{$}", h.postScanEvent)

	m.HandleFunc("GET /events/{id}/attendees/{$}", h.getEventAttendees)
	m.HandleFunc("POST /events/{id}/attendees/{$}", h.postEventAttendees)

	m.HandleFunc("GET /events/{id}/merge/{$}", h.getMergeEvent)
	m.HandleFunc("POST /events/{id}/merge/{$}", h.postMergeEvent)

	m.Handle("GET /", frontend.StaticHandler())

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, pattern := h.mux.Handler(r)
	h.logger.Debug(
		"Handling request",
		"method", r.Method,
		"path", r.URL.Path,
		"pattern", pattern)
	if pattern == "" {
		h.renderErrorWithCode(w, 404, fmt.Errorf("page not found"))
		return
	}
	handler.ServeHTTP(w, r)
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

type eventData struct {
	Event     sqldb.Event
	Attendees []sqldb.ListEventAttendeesRow
}

func (h *Handler) getEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	slog.Debug(
		"Getting event",
		"id", id)

	event, err := h.db.GetEvent(r.Context(), id)
	if err != nil {
		h.renderErrorWithCode(w, 500, fmt.Errorf("cannot get event: %w", err))
		return
	}

	attendees, err := h.db.ListEventAttendees(r.Context(), id)
	if err != nil {
		h.renderErrorWithCode(w, 500, fmt.Errorf("cannot get attendees: %w", err))
		return
	}

	h.tmpl.Execute(w, "event", eventData{
		Event:     event,
		Attendees: attendees,
	})
}

func (h *Handler) getScanEvent(w http.ResponseWriter, r *http.Request) {
	h.tmpl.Execute(w, "scan", nil)
}

func (h *Handler) postScanEvent(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) getEventAttendees(w http.ResponseWriter, r *http.Request) {
	h.tmpl.Execute(w, "event_attendees", nil)
}

func (h *Handler) postEventAttendees(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) getMergeEvent(w http.ResponseWriter, r *http.Request) {
	h.tmpl.Execute(w, "event_merge", nil)
}

func (h *Handler) postMergeEvent(w http.ResponseWriter, r *http.Request) {}
