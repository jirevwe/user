package util

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeJson(body io.Reader, o interface{}) error {
	return json.NewDecoder(body).Decode(o)
}

func EncodeJson(w http.ResponseWriter, err error) error {
	return json.NewEncoder(w).Encode(NewServiceErrResponse(err))
}

func EncodeJsonStatus(w http.ResponseWriter, message string, statusCode int, o interface{}) error {
	return json.NewEncoder(w).Encode(NewServerResponse(message, o, statusCode))
}
