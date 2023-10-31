package ioutils

import (
	"io/ioutil"
	"net/http"
)

//easyjson:json
type JSON struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

//easyjson:json
type ModelError struct {
	Message string `json:"message,omitempty"`
}

type ReadModel interface {
	UnmarshalJSON(data []byte) error
}

type WriteModel interface {
	MarshalJSON() ([]byte, error)
}

func Send(w http.ResponseWriter, respCode int, respBody WriteModel) {
	w.WriteHeader(respCode)
	_ = writeJSON(w, respBody)
}

func SendError(w http.ResponseWriter, respCode int, errorMsg string) {
	Send(w, respCode, ModelError{
		Message: errorMsg,
	})
}

func SendDefaultError(w http.ResponseWriter, respCode int) {
	Send(w, respCode, ModelError{
		Message: resolveErrorToString(respCode),
	})
}

func resolveErrorToString(respCode int) string {
	switch respCode {
	case http.StatusBadRequest:
		return "bad request"
	case http.StatusUnauthorized:
		return "no auth"
	case http.StatusForbidden:
		return "forbidden"
	case http.StatusNotFound:
		return "not found"
	case http.StatusConflict:
		return "conflict"
	case http.StatusInternalServerError:
		return "internal"
	default:
		return "unknown error"
	}
}

func SendWithoutBody(w http.ResponseWriter, respCode int) {
	w.WriteHeader(respCode)
}

func ReadJSON(r *http.Request, data ReadModel) error {
	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = data.UnmarshalJSON(byteReq)
	if err != nil {
		return err
	}

	return nil
}

func writeJSON(w http.ResponseWriter, data WriteModel) error {
	byteResp, err := data.MarshalJSON()
	if err != nil {
		return err
	}

	_, err = w.Write(byteResp)
	if err != nil {
		return err
	}

	return nil
}
