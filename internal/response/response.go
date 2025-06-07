package response

import (
	"fmt"
	"io"

	"github.com/ayrtonvitor/http/internal/headers"
)

const httpVersion = "HTTP/1.1"
const crlf = "\r\n"

func WriteStatusLine(w io.Writer, statusCode HttpStatusCode) error {
	var statusText string
	switch statusCode {
	case StatusCodeOk:
		statusText = string(StatusTextOk)
	case StatusCodeBadRequest:
		statusText = string(StatusTextBadRequest)
	case StatusCodeInternalServerError:
		statusText = string(StatusTextInternalServerError)
	default:
		statusText = ""
	}
	line := fmt.Sprintf("%s %d %s", httpVersion, statusCode, statusText)
	_, err := w.Write([]byte(line))
	return err
}
