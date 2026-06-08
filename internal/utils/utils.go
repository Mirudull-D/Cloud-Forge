package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
func WriteError(w http.ResponseWriter, status int, err error) {
	err2 := WriteJson(w, status, map[string]string{"error": err.Error()})
	if err2 != nil {
		log.Fatal(err2.Error())
	}
}
func ParseJson(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("no Payload Found")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}
