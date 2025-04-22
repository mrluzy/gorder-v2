package consts

const (
	ErrnoSuccess              = 0
	ErrnoUnknown              = 1
	ErrnoBindRequestError     = 1000
	ErrnoRequestValidateError = 1001
)

var ErrMsg = map[int]string{
	ErrnoSuccess:              "success",
	ErrnoUnknown:              "unknown error",
	ErrnoBindRequestError:     "bind request error",
	ErrnoRequestValidateError: "request validate error",
}
