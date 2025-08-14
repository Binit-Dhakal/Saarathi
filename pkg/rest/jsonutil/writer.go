package jsonutil

import (
	"encoding/json"
	"net/http"
)

type Writer struct {
}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) JSON(rw http.ResponseWriter, status int, data any) error {
	return w.JSONWithHeaders(rw, status, data, nil)
}

func (w *Writer) JSONWithHeaders(rw http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		rw.Header()[key] = value
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(js)

	return nil
}
