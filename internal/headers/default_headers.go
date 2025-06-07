package headers

import (
	"strconv"
	"strings"
)

const ContentLength = "content-length"
const Connection = "connection"
const ContentType = "content-type"

func GetDefaultHeaders(contentLen int) Headers {
	defaultHeaders := NewHeaders()
	defaultHeaders[GetAsCanonical(ContentLength)] = strconv.Itoa(contentLen)
	defaultHeaders[GetAsCanonical(Connection)] = "close"
	defaultHeaders[GetAsCanonical(ContentType)] = "text/plain"
	return defaultHeaders
}

func GetAsCanonical(key string) string {
	parts := strings.Split(key, "-")
	toJoin := []string{}
	for _, part := range parts {
		part = strings.ToLower(part)
		ini := part[0]
		if 'a' <= ini && ini <= 'z' {
			ini -= 'a' - 'A'
			toJoin = append(toJoin, string(ini)+part[1:])
		}
	}
	return strings.Join(toJoin, "-")
}
