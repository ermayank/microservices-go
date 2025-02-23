package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data, omitempty"`
}

// Function to Read any JSON
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes));

	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(data)

	if err != nil {
		return err
	}

	err = jsonDecoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("Body must have a single JSON Value")
	}

	return nil
}

// function to write any JSON
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	//Check for headers
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)

	if err != nil {
		return err
	}

	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)

}