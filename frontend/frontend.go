package frontend

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
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
