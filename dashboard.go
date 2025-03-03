package main

import (
  "html/template"
  "log"
  "os"
  "io"
  "fmt"
  "strings"
  "net/http"
  "sync"
  "gopkg.in/yaml.v3"
)

type Icon struct {
  name string
  html string
}

type Config struct {
  Title string `yaml:"title"`
  Layout Layout `yaml:"layout"`
  Style Style `yaml:"style"`
  LinkSections []LinkSection `yaml:"linkSections"`
}

type Layout struct {
  Sections int `yaml:"sections"`
  Width int `yaml:"width"`
  SectionPadding int `yaml:"sectionPadding"`
  CardPadding int `yaml:"cardPadding"`
}

type Style struct {
  Background string `yaml:"background"`
  SectionBackground string `yaml:"sectionBackground"`
  CardBackground string `yaml:"cardBackground"`
  CardHover string `yaml:"cardHover"`
  Text string `yaml:"text"`
  TextHover string `yaml:"textHover"`
  Accent string `yaml:"accent"`
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
var iconCache []Icon
var cacheMutex sync.Mutex

func main() {
  var err error
  config, err = loadConfig()
  
  if err != nil {
		log.Printf("error loading config: %v\n", err)
    os.Exit(1)
  }

  http.HandleFunc("/", indexHandler)
  err = http.ListenAndServe(":8080", nil)
  
  if err != nil {
    log.Printf("error starting listening on port :8080")
    os.Exit(1)
  }
}

func getIconUrl(icon string) string {
  if icon == "" {
    return "";
  }
  
  iconParts := strings.SplitN(icon, "-", 2)
  if len(iconParts) < 2 {
    return icon
  }

  switch iconParts[0] {
  case "hl":
    return "https://raw.githubusercontent.com/walkxcode/dashboard-icons/master/svg/" + iconParts[1] +".svg"
  
  case "fa":
    return "https://site-assets.fontawesome.com/releases/v5.15.4/svgs/regular" + iconParts[1] + ".svg"

  case "fas":
    return "https://site-assets.fontawesome.com/releases/v5.15.4/svgs/solid/" + iconParts[1] + ".svg"

  default:
    return icon
  }
}

func prefixSVGClasses(svg string, prefix string) string {
  replacer := strings.NewReplacer(
    ".st", "."+prefix+"-st",
    "class=\"st", "class=\""+prefix+"-st",
  )

  return replacer.Replace(svg)
}

func loadIcon(icon string) (template.HTML, error) {
  cacheMutex.Lock();
  for _, cachedIcon := range iconCache {
    if cachedIcon.name == icon {
      cacheMutex.Unlock()

      return template.HTML(cachedIcon.html), nil
    }
  }
  cacheMutex.Unlock()
 
  iconUrl := getIconUrl(icon)
  log.Printf("loading icon from: %v", iconUrl)
  
  res, err := http.Get(iconUrl)

  if err != nil {
    log.Printf("error making http request: %s\n", err)
    return "", fmt.Errorf("error making http request: %s\n", err)
  }
  
  if res.StatusCode != 200 {
    log.Printf("error getting icon: %v, code: %d", iconUrl, res.StatusCode)
    return "", fmt.Errorf("error getting icon: %v, code: %d\n", iconUrl, res.StatusCode)
  }

  defer res.Body.Close()
  body, err := io.ReadAll(res.Body)

  if err != nil {
    log.Printf("error getting icon from http body: %v", err)
    return "", fmt.Errorf("error getting icon http body %s\n", err)
  }

	prefix := strings.ReplaceAll(icon, "-", "_")
	svgWithScopedStyles := prefixSVGClasses(string(body), prefix)
  
  cacheMutex.Lock()
  iconCache = append(iconCache, Icon{name: icon, html: svgWithScopedStyles})
  cacheMutex.Unlock()

  log.Printf("loaded icon from http request: %v\n", icon)

  return template.HTML(svgWithScopedStyles), nil
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

  var wg sync.WaitGroup

  for _, linkSections := range cfg.LinkSections {
    for _, link := range linkSections.Links {
      if link.Icon != "" {
        wg.Add(1)
        go func(icon string) {
          defer wg.Done()
          
          _, _ = loadIcon(link.Icon)
        }(link.Icon)
      }
    }
  }

  wg.Wait()

  return cfg, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.New("index.html").Funcs(template.FuncMap{
    "loadIcon": loadIcon,
  })
  
  t, err := tmpl.ParseFiles("index.html")
  
  if err != nil {
    http.Error(w, "Error loading template", http.StatusInternalServerError)
    log.Printf("error loading template\n")
    
    return
  }

  err = t.Execute(w, config)

  if err != nil {
    http.Error(w, "Error rendering template", http.StatusInternalServerError)
    log.Printf("error rendering template\n")
    
    return
  }

  log.Printf("from %v served path %v\n", r.RemoteAddr, r.URL)
}

