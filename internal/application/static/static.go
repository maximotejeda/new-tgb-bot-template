package static

import (
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// RemoveAccent
// helps normalize names in db
// https://stackoverflow.com/questions/24588295/go-removing-accents-from-strings
func RemoveAccent(str string) string {
	if str == "" {
		return ""
	}
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ := transform.String(t, str)

	return s
}
