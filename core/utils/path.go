package utils

import (
	"regexp"
	"strings"
)

var colonPathParamRE = regexp.MustCompile(`:([a-zA-Z0-9_]+)`)
var templatedPathParamRE = regexp.MustCompile(`\{([a-zA-Z0-9_]+)\}`)

func NormalizeAPIPath(path string) string {
	normalized := strings.TrimSpace(path)
	if normalized == "" {
		return "/"
	}
	normalized = colonPathParamRE.ReplaceAllString(normalized, `{$1}`)
	if normalized != "/" {
		normalized = strings.TrimSuffix(normalized, "/")
		if normalized == "" {
			return "/"
		}
	}
	return normalized
}

func ExtractPathParamNames(path string) []string {
	matches := templatedPathParamRE.FindAllStringSubmatch(NormalizeAPIPath(path), -1)
	names := make([]string, 0, len(matches))
	for _, match := range matches {
		names = append(names, match[1])
	}
	return names
}
