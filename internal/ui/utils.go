package ui

import (
	"strings"
	"unicode"
)

func (m model) getCurrentServer() string {
	return m.table.SelectedRow()[0]
}

func (m model) humanizeError(err error) string {
	rawParts := strings.Split(err.Error(), "ssh: ")
	var cleanedParts []string

	for _, p := range rawParts {
		p = strings.TrimSpace(p)
		p = strings.TrimSuffix(p, ":")
		p = strings.TrimSpace(p)

		if p == "" {
			continue
		}

		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		cleanedParts = append(cleanedParts, string(runes))
	}

	result := strings.Join(cleanedParts, ". ")

	if len(result) > 0 && !strings.HasSuffix(result, ".") {
		result += "."
	}

	return result
}
