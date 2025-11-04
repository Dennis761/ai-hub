package prompttmpl

import (
	"net/url"
	"regexp"
	"sort"
	"strings"
)

// parse query params from apiMethod like "/v1/run?key=value&mode=test#frag"
func ParseQueryParamsFromAPIMethod(apiMethod string) map[string]string {
	res := make(map[string]string)

	apiMethod = strings.TrimSpace(apiMethod)
	if apiMethod == "" {
		return res
	}

	qIdx := strings.Index(apiMethod, "?")
	if qIdx == -1 {
		return res
	}

	raw := apiMethod[qIdx+1:]
	if hashIdx := strings.Index(raw, "#"); hashIdx >= 0 {
		raw = raw[:hashIdx]
	}

	values, err := url.ParseQuery(raw)
	if err != nil {
		return res
	}
	for k, v := range values {
		if len(v) > 0 {
			res[k] = v[len(v)-1]
		} else {
			res[k] = ""
		}
	}
	return res
}

// extract placeholders like {{ key }} (supports letters, digits, _, ., -)
func ExtractPlaceholders(text string) []string {
	if strings.TrimSpace(text) == "" {
		return []string{}
	}
	re := regexp.MustCompile(`{{\s*([\w.-]+)\s*}}`)
	matches := re.FindAllStringSubmatch(text, -1)

	seen := make(map[string]struct{}, len(matches))
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		key := strings.TrimSpace(m[1])
		if key == "" {
			continue
		}
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			out = append(out, key)
		}
	}
	return out
}

// replace all {{ key }} with values from params (longer keys first)
func FormatPrompt(promptText string, params map[string]string) string {
	filled := promptText
	if filled == "" || len(params) == 0 {
		return filled
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return len(keys[i]) > len(keys[j]) })

	for _, key := range keys {
		val := params[key]
		re := regexp.MustCompile(`{{\s*` + regexp.QuoteMeta(key) + `\s*}}`)
		filled = re.ReplaceAllString(filled, val)
	}
	return filled
}
