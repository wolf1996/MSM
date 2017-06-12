package view

import (
	"net/http"
	"encoding/json"
)

type msgObj interface{

}

type messageContainer struct {
	Code int `json:"code"`
	Msg msgObj `json:"message"`
}

func WriteMessage ( w *http.ResponseWriter, msg msgObj, res int) ErrorView {
	container := messageContainer{res, msg}
	if err := json.NewEncoder(*w).Encode(container); err != nil {
		return ErrorViewImpl{"Encoding error", 1}
	}
	return nil
}
