package server

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"

	"dev.acmcsuf.com/fullyhacks-qrms/sqldb"
	"github.com/skip2/go-qrcode"
	"github.com/sourcegraph/conc/iter"
)

func (h *Handler) listUsersPage(w http.ResponseWriter, r *http.Request) {
	users, err := h.db.ListUsers(r.Context())
	if err != nil {
		h.renderErrorWithCode(w, 500, fmt.Errorf("failed to list users: %w", err))
		return
	}

	h.tmpl.Execute(w, "user_list", map[string]any{
		"Users": users,
	})
}

func (h *Handler) addUserPage(w http.ResponseWriter, r *http.Request) {
	h.tmpl.Execute(w, "user_add", nil)
}

func (h *Handler) addUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	if name == "" || email == "" {
		h.renderErrorWithCode(w, 400, fmt.Errorf("malformed form"))
		return
	}

	user := sqldb.NewUser(name, email)

	if err := h.db.AddUser(r.Context(), sqldb.AddUserParams{
		Email: email,
		Code:  user.Code,
		Name:  name,
	}); err != nil {
		h.renderErrorWithCode(w, 500, fmt.Errorf("failed to create user: %w", err))
		return
	}

	http.Redirect(w, r, "/users#"+email, 303)
}

func (h *Handler) getAllUserQRs(w http.ResponseWriter, r *http.Request) {
	users, err := h.db.ListUsers(r.Context())
	if err != nil {
		h.renderErrorWithCode(w, 500, fmt.Errorf("failed to list users: %w", err))
		return
	}

	w.Header().Set("Trailers", "Content-Disposition, Content-Type")
	w.WriteHeader(200)

	// onError is a helper function to write the error in a special way.
	// It is required because the headers have already been written.
	onError := func(err error) {
		h.logger.Error(
			"Error while streaming QR codes as a zip file",
			"error", err,
			"users", len(users))

		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "\n\nInternal server error", 500)
	}

	zipFS := zip.NewWriter(w)

	mappingsFile, err := zipFS.Create("mappings.json")
	if err != nil {
		onError(fmt.Errorf("failed to create mappings file: %w", err))
		return
	}

	if err := renderQRMappingFile(users, mappingsFile); err != nil {
		onError(fmt.Errorf("failed to render mappings file: %w", err))
		return
	}

	type renderedPNG struct {
		Name string
		Data []byte
	}

	pngs, err := iter.MapErr(users, func(u *sqldb.User) (renderedPNG, error) {
		qrImage, err := renderQRImage(u.Code)
		if err != nil {
			return renderedPNG{}, fmt.Errorf("failed to render QR code: %w", err)
		}

		var b bytes.Buffer
		b.Grow(512)
		if err := renderImageToPNG(&b, qrImage); err != nil {
			return renderedPNG{}, fmt.Errorf("failed to render QR code to PNG: %w", err)
		}

		return renderedPNG{
			Name: u.Code + ".png",
			Data: b.Bytes(),
		}, nil
	})
	if err != nil {
		onError(fmt.Errorf("failed to render QR codes: %w", err))
		return
	}

	for _, png := range pngs {
		f, err := zipFS.Create(png.Name)
		if err != nil {
			onError(fmt.Errorf("failed to create QR file: %w", err))
			return
		}

		if _, err := f.Write(png.Data); err != nil {
			onError(fmt.Errorf("failed to write QR file: %w", err))
			return
		}
	}

	if err := zipFS.Close(); err != nil {
		onError(fmt.Errorf("failed to finalize zip file: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=qr_codes.zip")
}

func (h *Handler) getUserQRAsPNG(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")

	code, err := h.db.GetUserCode(r.Context(), email)
	if err != nil {
		h.renderErrorWithCode(w, 400, fmt.Errorf("failed to find user: %w", err))
		return
	}

	qrImage, err := renderQRImage(code)
	if err != nil {
		h.renderErrorWithCode(w, 500, fmt.Errorf("failed to render QR code: %w", err))
		return
	}

	w.Header().Set("Content-Type", "image/png")
	renderImageToPNG(w, qrImage)
}

// qrModuleSize is the size of each module (dot) in the QR code.
const qrModuleSize = 8

func renderQRImage(data string) (image.Image, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}

	return qr.Image(-qrModuleSize), nil
}

func renderImageToPNG(w io.Writer, img image.Image) error {
	return png.Encode(w, img)
}

func renderQRMappingFile(users []sqldb.User, w io.Writer) error {
	mappings := make(map[string]string, len(users))
	for _, u := range users {
		mappings[u.Code] = u.Email
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(mappings); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
