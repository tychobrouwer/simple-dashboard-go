package utils

import "regexp"

func PrefixSVGClasses(svg string, prefix string) string {
	re1 := regexp.MustCompile(`class="([a-z]+)`)
	re2 := regexp.MustCompile(`\.([a-z]+)`)
	re3 := regexp.MustCompile(`id="([a-z]+)"`)
	re4 := regexp.MustCompile(`([^:"])#([a-z]+)`)

	svg = re1.ReplaceAllString(svg, `class="`+prefix+`-$1`)
	svg = re2.ReplaceAllString(svg, `.`+prefix+`-$1`)
	svg = re3.ReplaceAllString(svg, `id="`+prefix+`-$1"`)
	svg = re4.ReplaceAllString(svg, `$1#`+prefix+`-$2`)

	return svg
}
