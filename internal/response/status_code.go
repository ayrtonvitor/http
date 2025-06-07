package response

type httpStatusCode int
type httpStatusText string

const (
	statusCodeOk                  httpStatusCode = 200
	statusCodeBadRequest          httpStatusCode = 400
	statusCodeInternalServerError httpStatusCode = 500
)

const (
	statusTextOk                  httpStatusText = "OK"
	statusTextBadRequest          httpStatusText = "Bad Request"
	statusTextInternalServerError httpStatusText = "Internal Server Error"
)
