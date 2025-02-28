package main

import (
  "html/template"
  "log"
  "os"
  "io"
  "fmt"
  "errors"
  "strings"
  "net/http"
  "gopkg.in/yaml.v3"
)

type Config struct {
  Layout Layout `yaml:"layout"`
  Style Style `yaml:"style"`
  LinkSections []LinkSection `yaml:"linkSections"`
}

type Layout struct {
  Sections int `yaml:"sections"`
  Width int `yaml:"width"`
}

type Style struct {
  AccentColor string `yaml:"accentColor"`
}

type LinkSection struct {
  Title string `yaml:"title"`
  Links []Links `yaml:"links"`
}

type Links struct {
  Title string `yaml:"title"`
  Link string `yaml:"link"`
  Icon string `yaml:"icon"`
}

var config Config

func main() {
  var err error
  config, err = loadConfig()
  
  if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
  }

  http.HandleFunc("/", indexHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func getIcon(icon string) (string, error) {
  if icon == "" {
    return "", nil;
  }
  
  if len(icon) < 3 {
    return "", errors.New("invalid icon string")
  }

  iconParts := strings.SplitN(icon, "-", 2)
  if len(iconParts) < 2 {
    return icon, nil
  }

  switch iconParts[0] {
  case "hl":
    return "https://raw.githubusercontent.com/walkxcode/dashboard-icons/master/svg/" + iconParts[1] +"/.svg", nil
  
  case "fa":
    return "https://site-assets.fontawesome.com/releases/v5.15.4/svgs/regular" + iconParts[1] + ".svg", nil

  case "fas":
    return "https://site-assets.fontawesome.com/releases/v5.15.4/svgs/solid/" + iconParts[1] + ".svg", nil

  default:
    return icon, nil
  }
}

func loadConfig() (Config, error) {
  filename := "config.yml"
  configFile, err := os.Open(filename)

  if err != nil {
		return Config{}, fmt.Errorf("error opening config file: %v\n", err)
  }

  byteValue, err := io.ReadAll(configFile) 
  
  if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %v\n", err)
  }
  
  defer configFile.Close()

  var cfg Config

  err = yaml.Unmarshal(byteValue, &cfg)

  if err != nil {
		return Config{}, fmt.Errorf("error parsing config file: %v\n", err)
  }

  return cfg, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("index.html").Funcs(template.FuncMap{
		"getIcon": getIcon,
	})
  
  t, err := tmpl.ParseFiles("index.html")
  
  if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
  }

  err = t.Execute(w, config)

  if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
  }

  fmt.Printf("from %v served path %v\n", r.RemoteAddr, r.URL)
}

