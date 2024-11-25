package consts

const (
	ErrnoSuccess               = 0
	ErrnoUnknownError          = 1
	ErrnoBindRequestsError     = 1000
	ErrnoRequestsValidateError = 1001
	// ErrnoOrder = 2xxx ?
)

var ErrMsg = map[int]string{
	ErrnoSuccess:               "success",
	ErrnoUnknownError:          "unknown error",
	ErrnoBindRequestsError:     "binding requests failed",
	ErrnoRequestsValidateError: "requests validation failed",
}
