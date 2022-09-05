package common

type HttpError interface {
	error
	HttpCode() int
}

type HttpErrorImpl struct {
	Code int
	err  error
}

func (err *HttpErrorImpl) Error() string {
	return err.err.Error()
}

func (err *HttpErrorImpl) HttpCode() int {
	return err.Code
}

func NewHttpError(code int, err error) HttpError {
	return &HttpErrorImpl{
		Code: code,
		err:  err,
	}
}
