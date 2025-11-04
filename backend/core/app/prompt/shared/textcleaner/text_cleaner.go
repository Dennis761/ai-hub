package textcleaner

import (
	"regexp"
	"strings"
)

// CleanText removes HTML, markdown, and normalizes whitespace.
func CleanText(input string) string {
	if strings.TrimSpace(input) == "" {
		return ""
	}

	reHTML := regexp.MustCompile(`<[^>]*>`)
	cleaned := reHTML.ReplaceAllString(input, "")

	reMarkdown := regexp.MustCompile(`[*_~` + "`" + `#>-]+`)
	cleaned = reMarkdown.ReplaceAllString(cleaned, "")

	reSpaces := regexp.MustCompile(`\s+`)
	cleaned = reSpaces.ReplaceAllString(cleaned, " ")

	return strings.TrimSpace(cleaned)
}
