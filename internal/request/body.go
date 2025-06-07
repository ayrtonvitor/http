package request

import (
	"errors"
)

type Body struct {
	content        []byte
	expectedLength int
}

var ErrReqBodyLongerThanReported = errors.New("Request body is longer than informed value")

func (b *Body) parse(data []byte) (int, bool, error) {
	if b.expectedLength == 0 {
		return 0, true, nil
	}
	b.content = append(b.content, data...)
	if len(b.content) > b.expectedLength {
		return 0, false, ErrReqBodyLongerThanReported
	}
	return len(data), len(b.content) == b.expectedLength, nil
}

func (b *Body) AsString() string {
	return string(b.content)
}
