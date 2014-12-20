package main

import (
	"regexp"
	"strings"
)

// Strip any special character.
var reStripFormat = regexp.MustCompile("\033\\[[\\d;]+m")

// Remove formatting from string
func StripFormat(s string) string {
	return strings.TrimSpace(reStripFormat.ReplaceAllString(s, ""))
}
