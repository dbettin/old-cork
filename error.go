package cork

type Error struct {
	error
	Status int
	//stacktrace
}

func NewError(err error, status int) *Error {
	return &Error{err, status}
}
