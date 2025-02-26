package main

import (
  "html/template"
  "log"
  "os"
  "io/ioutil"
  "fmt"
  "net/http"
  //"encoding/json"
  "gopkg.in/yaml.v3"
)

type Config struct {
  LinkSections []LinkSection `yaml:"linkSections"`
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

func main() {
  http.HandleFunc("/", indexHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadConfig() (*Config, error) {
  filename := "config.yml"
  configFile, errOpen := os.Open(filename)

  if errOpen != nil {
    fmt.Println("Error opening config file: %v", errOpen)

    return nil, errOpen
  }

  byteValue, errRead := ioutil.ReadAll(configFile) 
  
  if errRead != nil {
    fmt.Println("Error reading config file: %v", errRead)

    return nil, errRead
  }
  
  defer configFile.Close()

  var config Config

  errParse := yaml.Unmarshal(byteValue, &config)

  if errParse != nil {
    fmt.Println("Error parsing config file: %v", errParse)

    return nil, errParse
  }

  return &config, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
  config, err := loadConfig()
  if err != nil {
    //config = &Config{Links: []string{}}
  }


  for _, section := range config.LinkSections {
    fmt.Println("section: " + section.Title)

    for _, link := range section.Links {
      fmt.Println("link title: " + link.Title)
      fmt.Println("link link: " + link.Link)
      fmt.Println("link icon: " + link.Icon)
    }
  }

  t, _ := template.ParseFiles("index.html")
  t.Execute(w, config)
}

