package response

type HttpStatusCode int
type HttpStatusText string

const (
	StatusCodeOk                  HttpStatusCode = 200
	StatusCodeBadRequest          HttpStatusCode = 400
	StatusCodeInternalServerError HttpStatusCode = 500
)

const (
	StatusTextOk                  HttpStatusText = "OK"
	StatusTextBadRequest          HttpStatusText = "Bad Request"
	StatusTextInternalServerError HttpStatusText = "Internal Server Error"
)
