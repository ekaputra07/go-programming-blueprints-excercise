package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/objx"
)

func render(filename string) http.HandlerFunc {
	var tpl *template.Template
	var once sync.Once

	return func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() {
			log.Printf("Loading template %s\n", filename)
			tpl = template.Must(template.ParseFiles(filepath.Join("templates", filename)))
		})
		data := map[string]interface{}{
			"Host": r.Host,
		}
		if authCookie, err := r.Cookie("auth"); err == nil {
			data["UserData"] = objx.MustFromBase64(authCookie.Value)
		}
		tpl.Execute(w, data) // we pass Request as the data so we can access Request data in template
	}
}
