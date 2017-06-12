package view

type ErrorView interface {
	Error() string
	Id() int
}

type  ErrorViewImpl struct {
	Msg string
	Code int
}

func (e ErrorViewImpl)Error()string  {
	return e.Msg
}

func (e ErrorViewImpl)Id()int  {
	return e.Code
}
