package api

import "strings"

type Resource string

const ModelResource Resource = "models"

func resourceLink(r Resource, base string, args ...string) string {
	elems := []string{base, string(r)}
	elems = append(elems, args...)
	return strings.Join(elems, "/")
}
