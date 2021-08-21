package main

type MappedError struct {
	StatusCode int
	Message    string
	ErrorCode  string
}

func (e *MappedError) Error() string {
	return e.Message
}
