package config

import (
	"fmt"
	"icons"
	"io"
	"os"
	"sync"

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
}

type LinkSection struct {
	Title string  `yaml:"title"`
	Links []Links `yaml:"links"`
}

type Links struct {
	Title string `yaml:"title"`
	Link  string `yaml:"link"`
	Icon  string `yaml:"icon"`
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

					_, _, _ = icons.LoadIcon(link.Icon)
				}(link.Icon)
			}
		}
	}

	wg.Wait()

	return nil
}

func GetConfig() Config {
	return config
}
