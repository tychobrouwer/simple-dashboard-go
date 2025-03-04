package handlers

import (
	"config"
	"icons"
	"log"
	"net/http"
	"text/template"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("index.html").Funcs(template.FuncMap{
		"getIconSrc":  icons.GetIconSrc,
		"getIconHtml": icons.GetIconHtml,
	})

	t, err := tmpl.ParseFiles("index.html")

	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("error loading template\n")

		return
	}

	err = t.Execute(w, config.GetConfig())

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("error rendering template\n")

		return
	}

	log.Printf("from %v served path %v\n", r.RemoteAddr, r.URL)
}
