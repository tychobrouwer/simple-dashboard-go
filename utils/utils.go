package utils

import (
  "regexp"
  "os"
  "time"
  "strings"
)

func PrefixSVGClasses(svg string, prefix string) string {
	re1 := regexp.MustCompile(`class="([a-z]+)`)
	re2 := regexp.MustCompile(`([^[a-z0-9])\.([a-z]+)`)
	re3 := regexp.MustCompile(`id="([a-z]+)`)
	re4 := regexp.MustCompile(`([^:"])#([a-z]+)`)
  re5 := regexp.MustCompile(`xlink:href="#([a-z]+)`)

	svg = re1.ReplaceAllString(svg, `class="`       + prefix + `-$1`)
	svg = re2.ReplaceAllString(svg, `$1.`           + prefix + `-$2`)
	svg = re3.ReplaceAllString(svg, `id="`          + prefix + `-$1`)
	svg = re4.ReplaceAllString(svg, `$1#`           + prefix + `-$2`)
  svg = re5.ReplaceAllString(svg, `xlink:href="#` + prefix + `-$1`)

	return svg
}

func AddSVGViewBox(svg string) string {
  if strings.Contains(svg, "viewBox") {
    return svg
  }

  reWidth := regexp.MustCompile(`width="([a-z0-9]+)"`)
  attrWidth := reWidth.FindStringSubmatch(svg)
  reHeight := regexp.MustCompile(`height="([a-z0-9]+)"`)
  attrHeight := reHeight.FindStringSubmatch(svg)

  if len(attrWidth) == 2 && len(attrHeight) == 2 {
    viewbox := `<svg viewBox="0 0 ` + attrWidth[1] + ` ` + attrHeight[1] + `"`

    re := regexp.MustCompile(`<svg `)
    svg = re.ReplaceAllString(svg, viewbox)
  }

  return svg
}

func WatchFile(filePath string) error {
  initialStat, err := os.Stat(filePath)
  if err != nil {
    return err
  }

  for {
    stat, err := os.Stat(filePath)
    if err != nil {
      return err
    }

    if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
      break
    }
    
    time.Sleep(1 * time.Second)
  }
  
  return nil
}
