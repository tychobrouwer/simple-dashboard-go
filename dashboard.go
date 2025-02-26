package main

import (
  "html/template"
  "log"
  "os"
  "io/ioutil"
  "fmt"
  "net/http"
  "encoding/json"
)

type Config struct {
  LinkSections []LinkSection `json:"linkSections"`
}

type LinkSection struct {
  Title string `json:"title"`
  Links []Links `json:"links"`
}

type Links struct {
  Title string `json:"title"`
  Link string `json:"link"`
  Icon string `json:"icon"`
}

func main() {
  http.HandleFunc("/", indexHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadConfig() (*Config, error) {
  filename := "config.json"
  jsonFile, err := os.Open(filename)

  if err != nil {
      return nil, err
  }

  byteValue, _ := ioutil.ReadAll(jsonFile) 
  defer jsonFile.Close()

  var config Config

  json.Unmarshal(byteValue, &config)

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

