package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

const crlf = "\r\n"

var errMalformedHeader = errors.New("Malformed header")
var errEmptyValHeader = errors.New("Header does not contain a valid value")
var errInvalidFieldName = errors.New("Field name contains invalid characters")

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))

	if idx < 0 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}
	done, err = h.parseHeaderFromString(string(data[:idx]))
	if err != nil {
		return 0, false, err
	}

	return idx + 2, false, err
}

func (h Headers) parseHeaderFromString(str string) (done bool, err error) {
	idx := strings.Index(str, ":")

	if idx < 1 || str[idx-1] == ' ' {
		return false, errMalformedHeader
	}

	key := strings.TrimSpace(strings.ToLower(str[:idx]))
	if hasBadChar(key) {
		return false, errInvalidFieldName
	}

	val := strings.TrimSpace(str[idx+1:])
	if len(val) == 0 {
		return false, errEmptyValHeader
	}

	if oldVal, ok := h[key]; ok {
		val = oldVal + ", " + val
	}
	h[key] = val

	return true, nil
}

func hasBadChar(str string) bool {
	for _, r := range str {
		if _, ok := allowedChars[r]; !ok {
			return true
		}
	}
	return false
}

var allowedChars = func() map[rune]struct{} {
	const allowed = "abcdefghijklmnopqrstuvwxyz!#$%&'*+-.^_`|~"
	m := map[rune]struct{}{}
	for _, r := range allowed {
		m[r] = struct{}{}
	}
	return m
}()
