package util

import (
	"regexp"
	"strings"
)

// CreateSlug converts a string to a URL-friendly slug
// For example: "Hello World" -> "hello-world"
func CreateSlug(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace non-alphanumeric characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")

	// Remove leading and trailing hyphens
	s = strings.Trim(s, "-")

	return s
}
