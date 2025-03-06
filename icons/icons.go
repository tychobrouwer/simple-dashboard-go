package icons

import (
	"fmt"
	"html/template"
	"io"
  "encoding/base64"
	"log"
	"net/http"
	"net/url"
  "strings"
	"sync"
	"utils"
)

type Icon struct {
	Key  string
	Url  string
	Html string
	Src  string
}

var iconCache []Icon
var cacheMutex sync.Mutex

func getIconUrl(icon string, link string) (string, string) {
	if icon == "" {
		return "empty", ""
	}

  if icon == "favicon" {
    u, err := url.Parse(link)

    if err != nil {
      log.Printf("error parsing url for favicon: %v\n", err)
      return "favicon", link + "/favicon.ico"
    }

    return "favicon", u.Scheme + "://" + u.Host + "/favicon.ico"
  }

  if strings.HasSuffix(icon, ".svg") {
    return "svg", icon
  }

	iconParts := strings.SplitN(icon, "-", 2)
	if len(iconParts) < 2 {
		return "url", icon
	}

	switch iconParts[0] {
	case "hl":
		return "hl", "https://raw.githubusercontent.com/homarr-labs/dashboard-icons/master/svg/" + iconParts[1] + ".svg"

	case "fa":
		return "fa", "https://site-assets.fontawesome.com/releases/v5.15.4/svgs/regular" + iconParts[1] + ".svg"

	case "fas":
		return "fas", "https://site-assets.fontawesome.com/releases/v5.15.4/svgs/solid/" + iconParts[1] + ".svg"

	default:
		return "url", icon
	}
}

func GetIconSrc(icon string, url string) (string, error) {
  iconSrc, _, err := getCachedIcon(icon, url)

  if err == nil {
		return iconSrc, nil
	}

	iconSrc, _, err = LoadIcon(icon, url)

	return iconSrc, err
}

func GetIconHtml(icon string, url string) (template.HTML, error) {
  _, iconHtml, err := getCachedIcon(icon, url)
  
  if err == nil {
		return iconHtml, nil
	}

	_, iconHtml, err = LoadIcon(icon, url)

	return iconHtml, err
}

func LoadIcon(icon string, url string) (string, template.HTML, error) {
	iconSrc, iconHtml, err := getCachedIcon(icon, url)
	if err == nil {
		return iconSrc, iconHtml, nil
	}

	iconSrc, iconUrl := getIconUrl(icon, url)

  if iconSrc == "favicon" || iconSrc == "url" {
    imageHtml := "<img alt=\"" + icon + "\" src=\"" + iconUrl + "\">"

    saveCachedIcon(icon, url, imageHtml, iconSrc)

	  log.Printf("loaded image tag: %v\n", iconUrl)
    
    return iconSrc, template.HTML(imageHtml), nil
  }

	res, err := http.Get(iconUrl)

	if err != nil {
		log.Printf("error making http request: %s\n", err)
		return "empty", "", fmt.Errorf("error making http request: %s", err)
	}

	if res.StatusCode != 200 {
		log.Printf("error getting icon: %v, code: %d", iconUrl, res.StatusCode)
		return "empty", "", fmt.Errorf("error getting icon: %v, code: %d", iconUrl, res.StatusCode)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Printf("error getting icon from http body: %v", err)
		return "empty", "", fmt.Errorf("error getting icon http body %s", err)
	}

	prefix := strings.ReplaceAll(icon, "-", "_")
	
  svg := utils.PrefixSVGClasses(string(body), prefix)
  svg = utils.AddSVGViewBox(svg)

	saveCachedIcon(icon, iconUrl, svg, iconSrc)

	log.Printf("loaded svg icon: %v\n", icon)

	return iconSrc, template.HTML(svg), nil
}

func saveCachedIcon(icon string, url string, html string, src string) {
  if icon == "favicon" {
    icon = url
  }
  
  key := base64.StdEncoding.EncodeToString([]byte(icon))

  cacheMutex.Lock()
	iconCache = append(iconCache, Icon{Key: key, Html: html, Src: src})
	cacheMutex.Unlock()
}

func getCachedIcon(icon string, url string) (string, template.HTML, error) {
  if icon == "favicon" {
    icon = url
  }

  key := base64.StdEncoding.EncodeToString([]byte(icon))

	cacheMutex.Lock()
	for _, cachedIcon := range iconCache {
		if cachedIcon.Key == key {
			cacheMutex.Unlock()

			return cachedIcon.Src, template.HTML(cachedIcon.Html), nil
		}
	}
	cacheMutex.Unlock()

	return "", "", fmt.Errorf("icon not found in cache")
}
