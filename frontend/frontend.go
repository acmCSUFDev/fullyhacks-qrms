package frontend

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"libdb.so/tmplutil"
)

//go:embed *
var embedFS embed.FS

// NewTemplater returns a new templater with the given filesystem.
func NewTemplater() *tmplutil.Templater {
	t := &tmplutil.Templater{
		FileSystem: embedFS,
		Includes: map[string]string{
			"head":       "components/head.html",
			"header":     "components/header.html",
			"components": "components/components.html",
		},
		Functions: template.FuncMap{
			"rfc3339": func() string { return time.RFC3339 },
			"rfc822":  func() string { return time.RFC822 },
			"shortUUID": func(uuid string) string {
				part, _, _ := strings.Cut(uuid, "-")
				return part
			},
			"plural": func(n int, singular, plural string) string {
				if n == 1 {
					return fmt.Sprintf("%d %s", n, singular)
				}
				return fmt.Sprintf("%d %s", n, plural)
			},
		},
	}
	if err := t.Preregister("pages"); err != nil {
		panic(err)
	}
	return t
}

// StaticHandler returns a handler for serving static files.
func StaticHandler() http.Handler {
	fs_, _ := fs.Sub(embedFS, "static")
	if fs_ == nil {
		panic("static files not found")
	}
	return http.FileServer(http.FS(fs_))
}
