package handlers

import (
	"config"
	"encoding/json"
	"icons"
	"log"
	"net/http"
	"slices"
	"sync"
	"text/template"
)

var statusMutex sync.Mutex

var acceptedStatusCodes = []int{
	http.StatusOK,
	http.StatusUnauthorized,
	http.StatusForbidden,
}

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

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetConfig()

	result := make(map[string]bool)

	var wg sync.WaitGroup

	for _, linkSections := range cfg.LinkSections {
		for _, link := range linkSections.Links {
			if link.Status {
				wg.Add(1)
				go func(link config.Link) {
					defer wg.Done()

					status := getStatus(link)

					statusMutex.Lock()
					result[link.Title] = status
					statusMutex.Unlock()
				}(link)
			}
		}
	}

	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status": "ok",
		"result": result,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	log.Printf("from %v served path %v\n", r.RemoteAddr, r.URL)
}

func getStatus(link config.Link) bool {
	url := link.StatusUrl
	if link.StatusUrl == "" {
		url = link.Link
	}

	resp, err := http.Get(url)

	if err != nil {
		log.Printf("Error checking status for %s: %v\n", link.Title, err)

		return false
	}

	defer resp.Body.Close()

	return slices.Contains(acceptedStatusCodes, resp.StatusCode)
}
