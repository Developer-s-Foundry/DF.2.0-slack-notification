package utils

import (
	"encoding/json"
	"net/http"
)

func WriteToJson(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

func ReadDataFromJson(r *http.Request, reader interface{}) error {
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&reader); err != nil {
		return err
	}
	return nil
}
