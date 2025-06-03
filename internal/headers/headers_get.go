package headers

import "strings"

func (h Headers) Get(key string) string {
	key = strings.ToLower(key)
	return h[key]
}
