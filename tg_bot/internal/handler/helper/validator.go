package helper

import "regexp"

func IsValidShortLink(s string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z-]+$`, s)
	return matched
}
