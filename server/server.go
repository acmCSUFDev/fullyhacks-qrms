package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"dev.acmcsuf.com/fullyhacks-qrms/frontend"
	"dev.acmcsuf.com/fullyhacks-qrms/sqldb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"libdb.so/hrt"
	"libdb.so/tmplutil"
)

// Handler is the type for handling API requests.
type Handler struct {
	mux       *chi.Mux
	db        *sqldb.Database
	tmpl      *tmplutil.Templater
	logger    *slog.Logger
	rootToken string
}

type Args struct {
	Database  *sqldb.Database
	Logger    *slog.Logger
	RootToken string
}

// NewHandler creates a new Handler.
func NewHandler(args Args) *Handler {
	h := &Handler{
		mux:       chi.NewMux(),
		db:        args.Database,
		tmpl:      frontend.NewTemplater(),
		logger:    args.Logger,
		rootToken: args.RootToken,
	}

	h.tmpl.OnRenderFail = func(sub *tmplutil.Subtemplate, w io.Writer, err error) {
		h.renderError(w, err)
	}

	m := h.mux

	m.Use(middleware.CleanPath)
	m.Use(hrt.Use(hrt.Opts{
		Encoder:     hrt.JSONEncoder,
		ErrorWriter: hrt.JSONErrorWriter("error"),
	}))

	m.NotFound(func(w http.ResponseWriter, r *http.Request) {
		h.renderErrorWithCode(w, 404, fmt.Errorf("page not found"))
	})

	m.Group(func(m chi.Router) {
		// m.Use(slogchi.NewWithConfig(h.logger, slogchi.Config{
		// 	DefaultLevel:     slog.LevelDebug,
		// 	ClientErrorLevel: slog.LevelDebug,
		// 	ServerErrorLevel: slog.LevelWarn,
		// }))

		m.Get("/auth/{token}", h.getAuth)

		m.Group(func(m chi.Router) {
			m.Use(h.useAuth)
			m.Use(h.requireAuth)

			m.Get("/", h.getIndex)

			m.Route("/users", func(m chi.Router) {
				m.Get("/", h.listUsersPage)
				m.Get("/{email}/qr.png", h.getUserQRAsPNG)

				m.Get("/add", h.addUserPage)
				m.Post("/add", h.addUser)

				// This endpoint is super heavy, so we throttle it really hard.
				m.With(middleware.Throttle(1)).Get("/qr_codes.zip", h.getAllUserQRs)
			})

			m.Route("/events", func(m chi.Router) {
				m.Get("/new", h.getNewEvent)
				m.Post("/new", h.postNewEvent)

				m.Route("/{id}", func(m chi.Router) {
					m.Get("/", h.getEvent)
					m.Get("/scan", h.getScanEvent)

					m.Get("/attendees", h.getEventAttendees)
					m.Post("/attendees", hrt.Wrap(h.addEventAttendee))
					m.Post("/attendees/{email}/delete", h.removeEventAttendee)

					m.Get("/merge", h.getMergeEvent)
					m.Post("/merge", h.postMergeEvent)
				})
			})
		})
	})

	m.Mount("/", frontend.StaticHandler())

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) renderErrorWithCode(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
