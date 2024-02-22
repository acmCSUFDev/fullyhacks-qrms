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

type addEventAttendeeResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handler) addEventAttendee(ctx context.Context, r addEventAttendeeRequest) (addEventAttendeeResponse, error) {
	eventUUID := hrt.RequestFromContext(ctx).PathValue("id")

	var resp addEventAttendeeResponse

	u, err := h.db.GetUserFromCode(ctx, r.UserCode)
	if err != nil {
		if sqldb.IsNotFound(err) {
			return resp, hrt.NewHTTPError(400, "user not found")
		}
		return resp, fmt.Errorf("cannot get user: %w", err)
	}

	if err := h.db.AddAttendee(ctx, sqldb.AddAttendeeParams{
		EventUUID: eventUUID,
		UserCode:  r.UserCode,
	}); err != nil {
		if sqldb.IsUniqueConstraintError(err) {
			return resp, hrt.NewHTTPError(400, "user already attends event")
		}
		return resp, fmt.Errorf("cannot add attendee: %w", err)
	}

	h.logger.Info(
		"Added attendee",
		"event", eventUUID,
		"user_code", r.UserCode,
		"user_name", u.Name,
		"user_email", u.Email)

	return addEventAttendeeResponse{
		Name:  u.Name,
		Email: u.Email,
	}, nil
}

func (h *Handler) removeEventAttendee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	email := r.PathValue("email")

	_, err := h.db.RemoveAttendee(r.Context(), sqldb.RemoveAttendeeParams{
		EventUUID: id,
		Email:     email,
	})
	if err != nil {
		h.renderErrorWithCode(w, 400, fmt.Errorf("cannot remove attendee: %w", err))
		return
	}

	http.Redirect(w, r, "/events/"+id+"/attendees", http.StatusSeeOther)
}

func (h *Handler) getMergeEvent(w http.ResponseWriter, r *http.Request) {
	h.tmpl.Execute(w, "event_merge", nil)
}

func (h *Handler) postMergeEvent(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
