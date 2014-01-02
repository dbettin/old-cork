package cork

type ErrorHandler interface {
	Handle(*Response, *Request, Error)
}

type Error struct {
	error
	status int
	//stacktrace
}
