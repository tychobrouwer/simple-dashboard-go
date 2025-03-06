package config

import (
	"fmt"
	"icons"
	"io"
	"os"
	"sync"
  "utils"
  "log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Title        string        `yaml:"title"`
	Layout       Layout        `yaml:"layout"`
	Style        Style         `yaml:"style"`
	LinkSections []LinkSection `yaml:"linkSections"`
}

type Layout struct {
	Sections       int `yaml:"sections"`
	Width          int `yaml:"width"`
	SectionPadding int `yaml:"sectionPadding"`
	CardPadding    int `yaml:"cardPadding"`
}

type Style struct {
	Background        string `yaml:"background"`
	SectionBackground string `yaml:"sectionBackground"`
	CardBackground    string `yaml:"cardBackground"`
	CardHover         string `yaml:"cardHover"`
	Text              string `yaml:"text"`
	TextHover         string `yaml:"textHover"`
	Accent            string `yaml:"accent"`
  StatusOnline      string `yaml:"statusOnline"`
  StatusOffline     string `yaml:"statusOffline"`
}

type LinkSection struct {
	Title string `yaml:"title"`
	Links []Link `yaml:"links"`
}

type Link struct {
	Title     string `yaml:"title"`
	Url       string `yaml:"url"`
	Icon      string `yaml:"icon"`
	Status    bool   `yaml:"status"`
	StatusUrl string `yaml:"statusUrl"`
}

var config Config

func LoadConfig() error {
	filename := "config.yml"
	configFile, err := os.Open(filename)

	if err != nil {
		return fmt.Errorf("error opening config file: %v", err)
	}

	byteValue, err := io.ReadAll(configFile)

	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	defer configFile.Close()

	err = yaml.Unmarshal(byteValue, &config)

	if err != nil {
		return fmt.Errorf("error parsing config file: %v", err)
	}

	var wg sync.WaitGroup

	for _, linkSections := range config.LinkSections {
		for _, link := range linkSections.Links {
			if link.Icon != "" {
				wg.Add(1)
				go func(icon string) {
					defer wg.Done()

					_, _, _ = icons.LoadIcon(link.Icon, link.Url)
				}(link.Icon)
			}
		}
	}

	wg.Wait()

	return nil
}

func WatchConfig() {
  configWatchChan := make(chan bool)
  
  go func(configWatchChan chan bool) {
    defer func() {
      configWatchChan <- true
    }()

    err := utils.WatchFile("config.yml")

    if err != nil {
      log.Printf("error watching config: %v\n", err)

      return
    }

    fmt.Println("config file changed")
    
    err = LoadConfig()
    if err != nil {
      log.Printf("error updating config: %v\n", err)
    }

    WatchConfig()
  }(configWatchChan)
}

func GetConfig() Config {
	return config
}
