package service

import "regexp"

const (
	likedYouDefaultPageSize = 5
	likedYouMaxPageSize     = 100
)

var uuidRegex = regexp.MustCompile(`^[a-fA-F0-9-]{36}$`)
