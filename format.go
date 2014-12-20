package main

import "regexp"

// Strip any special character.
var reStrip = regexp.MustCompile("\\0.*m")

// DeColorString removes all color from the given string
func Strip(s []byte) []byte {
	return reStrip.ReplaceAll(s, []byte{})
}
