package main

import (
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	template *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.template = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.template.Execute(w, r) // we pass Request as the data so we can access Request data in template
}

func renderTemplate(filename string) *templateHandler {
	return &templateHandler{filename: filename}
}
