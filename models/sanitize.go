package models

import "github.com/microcosm-cc/bluemonday"

// strip it all, no html
func Strip(s string) string {
	p := bluemonday.StrictPolicy()
	return p.Sanitize(s)
}
