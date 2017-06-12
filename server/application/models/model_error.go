package models

type ErrorModel interface {
	Error() string
	Id() int
}

type  ErrorModelImpl struct {
	Msg string
	Code int
}

func (e ErrorModelImpl)Error()string  {
	return e.Msg
}

func (e ErrorModelImpl)Id() int  {
	return e.Code
}