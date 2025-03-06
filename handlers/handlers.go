package handlers

import (
  "assets"
	"config"
	"encoding/json"
	"icons"
	"log"
	"net/http"
	"slices"
  "compress/gzip"
  "strings"
  "sync"
	"text/template"
)

var statusMutex sync.Mutex

var acceptedStatusCodes = []int{
	http.StatusOK,
	http.StatusUnauthorized,
	http.StatusForbidden,
}

func RobotsTxtHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	robotsTxt := "User-agent: *\nDisallow: /"

	_, err := w.Write([]byte(robotsTxt))

  if err != nil {
		log.Printf("Error writing robots.txt: %v", err)
	}
}

type GzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (gzw *GzipResponseWriter) Write(b []byte) (int, error) {
	return gzw.Writer.Write(b)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		w = &GzipResponseWriter{ResponseWriter: w, Writer: gz}
	}

  tmpl := template.New("index.html").Funcs(template.FuncMap{
		"getIconSrc":  icons.GetIconSrc,
		"getIconHtml": icons.GetIconHtml,
	})

  tmpl, err := tmpl.ParseFS(assets.PublicFS, "public/index.html")
	
  if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("error loading template\n")

		return
	}

  w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, config.GetConfig())

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
		url = link.Url	}

	resp, err := http.Get(url)

	if err != nil {
		log.Printf("Error checking status for %s: %v\n", link.Title, err)

		return false
	}

	defer resp.Body.Close()

	return slices.Contains(acceptedStatusCodes, resp.StatusCode)
}
