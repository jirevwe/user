package util

import (
	"encoding/json"
	"io"
)

func DecodeJson(body io.Reader, o interface{}) error {
	return json.NewDecoder(body).Decode(o)
}

func EncodeJson(w io.Writer, err error) error {
	return json.NewEncoder(w).Encode(NewServiceErrResponse(err))
}

func EncodeJsonStatus(w io.Writer, message string, statusCode int, o interface{}) error {
	return json.NewEncoder(w).Encode(NewServerResponse(message, o, statusCode))
}
