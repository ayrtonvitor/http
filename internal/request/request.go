package request

import (
	"bytes"
	"errors"
	"io"
	"slices"
	"strings"

	"github.com/ayrtonvitor/http/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        Body
	state       reqState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type reqState = int

const (
	reqStateInitialized reqState = iota
	reqStateDone
	reqStateParsingHeaders
)

const crlf = "\r\n"
const bufSize = 8

var errUnknownReqState = errors.New("Unknown state")
var errProcReqInDoneState = errors.New("Trying to read data in done state")

var ErrMalformedReq = errors.New("Request is malformed")
var ErrUnknownVerb = errors.New("HTTP verb is invalid")
var ErrUnsupportedVersion = errors.New("HTTP version is not supported")

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{
		state:   reqStateInitialized,
		Headers: headers.NewHeaders(),
	}
	buf := make([]byte, bufSize)
	nRead := 0
	eofHit := false

	for req.state != reqStateDone {
		if nRead >= len(buf) {
			buf = append(buf, make([]byte, len(buf))...)
		}

		nReadNew, err := reader.Read(buf[nRead:])
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
		if errors.Is(err, io.EOF) && req.state != reqStateDone {
			if eofHit {
				return nil, ErrMalformedReq
			}
			eofHit = true
		}
		nRead += nReadNew

		nParsed, err := req.parse(buf[:nRead], eofHit)
		if err != nil {
			return nil, err
		}

		if nParsed > 0 {
			copy(buf, buf[nParsed:])
			nRead -= nParsed
		}
	}

	return req, nil
}

func (r *Request) parse(data []byte, reqEnd bool) (int, error) {
	switch r.state {
	case reqStateInitialized:
		n, rlDone, err := r.parseInitReq(data)
		return r.handleInnerParseReturn(n, err, reqEnd, rlDone, reqStateParsingHeaders)
	case reqStateParsingHeaders:
		n, hDone, err := r.Headers.Parse(data)
		return r.handleInnerParseReturn(n, err, reqEnd, hDone, reqStateDone)
	case reqStateDone:
		return 0, errProcReqInDoneState
	default:
		return 0, errUnknownReqState
	}
}

func (r *Request) handleInnerParseReturn(n int, err error, reqEnd, done bool, nextState reqState) (int, error) {
	if err != nil {
		return 0, err
	}
	if reqEnd {
		if !done {
			return 0, ErrMalformedReq
		}
		r.state = reqStateDone
		return n, nil
	}
	if done {
		r.state = nextState
		return n, nil
	}
	return n, nil
}

func (r *Request) parseInitReq(data []byte) (int, bool, error) {
	reqLine, n, err := parseRequestLine(data)
	if err != nil {
		return 0, false, err
	}

	if n == 0 {
		return 0, false, nil
	}

	r.RequestLine = *reqLine
	return n, true, nil
}

func parseRequestLine(rawReq []byte) (*RequestLine, int, error) {
	idx := bytes.Index(rawReq, []byte(crlf))
	if idx < 0 {
		return nil, 0, nil
	}

	rawReqLine := string(rawReq[:idx])
	reqLine, err := requestLineFromString(rawReqLine)
	if err != nil {
		return nil, 0, err
	}

	return reqLine, idx + len(crlf), nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	reqLineParts := strings.Split(str, " ")
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

	requestLine := &RequestLine{
		HttpVersion:   version,
		RequestTarget: reqLineParts[1],
		Method:        reqLineParts[0],
	}

	return requestLine, nil
}
