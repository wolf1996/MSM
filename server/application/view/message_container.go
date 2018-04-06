package view

import (
	"encoding/json"
	"net/http"
	"github.com/wolf1996/MSM/server/application/error_codes"
)

type msgObj interface {
}

type messageContainer struct {
	Code int    `json:"code"`
	Msg  msgObj `json:"msg"`
}

func WriteMessage(w *http.ResponseWriter, msg msgObj, res int) ErrorView {
	container := messageContainer{res, msg}
	if err := json.NewEncoder(*w).Encode(container); err != nil {
		return ErrorViewImpl{"Encoding error", error_codes.ENCODING_ERROR}
	}
	return nil
}
