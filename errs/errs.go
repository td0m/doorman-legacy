package errs


type Error struct {
	Code int
	Err string
}

func (e Error) Error() string {
	return e.Err
}

func New(code int, s string) Error {
	return Error{Code: code, Err: s}
}
