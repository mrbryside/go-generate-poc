package mystr

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func CamelCase(s string) string {
	s = regexp.MustCompile("[^a-zA-Z0-9_ ]+").ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "_", " ")
	s = cases.Title(language.AmericanEnglish, cases.NoLower).String(s)
	s = strings.ReplaceAll(s, " ", "")
	if len(s) > 0 {
		s = strings.ToLower(s[:1]) + s[1:]
	}

	return s
}

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}
