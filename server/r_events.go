package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"dev.acmcsuf.com/fullyhacks-qrms/sqldb"
	"github.com/go-chi/chi/v5"
	"libdb.so/hrt"
)

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

func (h *Handler) getEventAttendees(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	event, err := h.db.GetEvent(r.Context(), id)
	if err != nil {
		h.renderErrorWithCode(w, 400, fmt.Errorf("cannot get event: %w", err))
		return
	}

	attendees, err := h.db.ListEventAttendees(r.Context(), id)
	if err != nil {
		h.renderErrorWithCode(w, 400, fmt.Errorf("cannot get attendees: %w", err))
		return
	}

	h.tmpl.Execute(w, "event_attendees", struct {
		sqldb.GetEventRow
		Attendees []sqldb.ListEventAttendeesRow
	}{
		GetEventRow: event,
		Attendees:   attendees,
	})
}

type addEventAttendeeRequest struct {
	UserCode string `json:"user_code"`
}

func (h *Handler) addEventAttendee(ctx context.Context, r addEventAttendeeRequest) (hrt.None, error) {
	eventUUID := hrt.RequestFromContext(ctx).PathValue("id")
	if err := h.db.AddAttendee(ctx, sqldb.AddAttendeeParams{
		EventUUID: eventUUID,
		UserCode:  r.UserCode,
	}); err != nil {
		return hrt.Empty, fmt.Errorf("cannot add attendee: %w", err)
	}
	return hrt.Empty, nil
}

func (h *Handler) getMergeEvent(w http.ResponseWriter, r *http.Request) {
	h.tmpl.Execute(w, "event_merge", nil)
}

func (h *Handler) postMergeEvent(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
