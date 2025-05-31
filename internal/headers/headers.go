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

var errMalformedHeader = errors.New("Malformed heder")
var errEmptyValHeader = errors.New("Header does not contain a valid value")

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

	key := strings.TrimSpace(str[:idx])
	val := strings.TrimSpace(str[idx+1:])
	if len(val) == 0 {
		return false, errEmptyValHeader
	}

	h[key] = val

	return true, nil
}
