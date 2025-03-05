package icons

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"utils"
)

type Icon struct {
	Name string
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
    return "favicon", link + "/favicon.ico"
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
  var iconSrc string
  var err error

  if icon == "favicon" {
    iconSrc, _, err = getCachedIcon(url)
  } else {
    iconSrc, _, err = getCachedIcon(icon)
  }

  if err == nil {
		return iconSrc, nil
	}

	iconSrc, _, err = LoadIcon(icon, url)

	return iconSrc, err
}

func GetIconHtml(icon string, url string) (template.HTML, error) {
  var iconHtml template.HTML
  var err error

  if icon == "favicon" {
    _, iconHtml, err = getCachedIcon(url)
  } else {
    _, iconHtml, err = getCachedIcon(icon)
  }
  
  if err == nil {
		return iconHtml, nil
	}

	_, iconHtml, err = LoadIcon(icon, url)

	return iconHtml, err
}

func LoadIcon(icon string, url string) (string, template.HTML, error) {
	iconSrc, iconHtml, err := getCachedIcon(icon)
	if err == nil {
		return iconSrc, iconHtml, nil
	}

	iconSrc, iconUrl := getIconUrl(icon, url)
	log.Printf("loading icon from: %v", iconUrl)

  if iconSrc == "favicon" {
    log.Println("loading favicon icon")

    faviconHtml := "<img alt=\"" + icon + "\" src=\"" + iconUrl + "\">"

	  cacheMutex.Lock()
	  iconCache = append(iconCache, Icon{Name: url, Html: faviconHtml, Src: iconSrc})
	  cacheMutex.Unlock()
    
    return "favicon", template.HTML(faviconHtml), nil
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
	svgWithScopedStyles := utils.PrefixSVGClasses(string(body), prefix)

	cacheMutex.Lock()
	iconCache = append(iconCache, Icon{Name: icon, Html: svgWithScopedStyles, Src: iconSrc})
	cacheMutex.Unlock()

	log.Printf("loaded icon from http request: %v\n", icon)

	return iconSrc, template.HTML(svgWithScopedStyles), nil
}

func getCachedIcon(icon string) (string, template.HTML, error) {
	cacheMutex.Lock()
	for _, cachedIcon := range iconCache {
		if cachedIcon.Name == icon {
			cacheMutex.Unlock()

			return cachedIcon.Src, template.HTML(cachedIcon.Html), nil
		}
	}
	cacheMutex.Unlock()

	return "", "", fmt.Errorf("icon not found in cache")
}
