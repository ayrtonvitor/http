package response

import (
	"fmt"
	"io"
)

const httpVersion = "HTTP/1.1"

func WriteStatusLine(w io.Writer, statusCode httpStatusCode) error {
	var statusText string
	switch statusCode {
	case statusCodeOk:
		statusText = string(statusTextOk)
	case statusCodeBadRequest:
		statusText = string(statusTextBadRequest)
	case statusCodeInternalServerError:
		statusText = string(statusTextInternalServerError)
	default:
		statusText = ""
	}
	line := fmt.Sprintf("%s %d %s", httpVersion, statusCode, statusText)
	_, err := w.Write([]byte(line))
	return err
}
