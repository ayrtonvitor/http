package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

var ErrMalformedReq = errors.New("Request is malformed")
var ErrUnknownVerb = errors.New("HTTP verb is invalid")
var ErrUnsupportedVersion = errors.New("HTTP version is not supported")

func RequestFromReader(reader io.Reader) (*Request, error) {
	rawReq, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("Error reading request: %w\n", err)
	}
	reqLine, err := parseRequestLine(string(rawReq))
	if err != nil {
		return nil, fmt.Errorf("Error parsing request line: %w\n", err)
	}
	return &Request{
		RequestLine: *reqLine,
	}, nil
}

func parseRequestLine(rawReq string) (*RequestLine, error) {
	parts := strings.Split(rawReq, crlf)
	if len(parts) < 3 {
		return nil, ErrMalformedReq
	}
	reqLineParts := strings.Split(parts[0], " ")
	if len(reqLineParts) != 3 {
		return nil, ErrMalformedReq
	}

	verbs := []string{"GET", "POST"}
	if !slices.Contains(verbs, reqLineParts[0]) {
		return nil, ErrUnknownVerb
	}
	if reqLineParts[2] != "HTTP/1.1" {
		return nil, ErrUnsupportedVersion
	}
	version := "1.1"

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: reqLineParts[1],
		Method:        reqLineParts[0],
	}, nil
}
