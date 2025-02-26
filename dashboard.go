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
  configFile, err := os.Open(filename)

  if err != nil {
      return nil, err
  }

  byteValue, _ := ioutil.ReadAll(configFile) 
  defer configFile.Close()

  var config Config

  yaml.Unmarshal(byteValue, &config)

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

