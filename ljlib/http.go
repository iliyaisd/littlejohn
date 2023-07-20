package ljlib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrorHTTP struct {
	Message string `json:"message"`
}

func ResponseHTTP(w http.ResponseWriter, code int, payload interface{}) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		payloadBytes = []byte(fmt.Sprintf(`{"message":"cannot marshal HTTP payload: %s"}`, err))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(code)
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(payloadBytes)
	if err != nil {
		log.Printf("cannot write HTTP response: %s", err)
	}
}

func ResponseHTTPError(w http.ResponseWriter, message string) {
	ResponseHTTP(w, http.StatusInternalServerError, ErrorHTTP{Message: message})
}

func ResponseHTTPBadRequest(w http.ResponseWriter, message string) {
	ResponseHTTP(w, http.StatusBadRequest, ErrorHTTP{Message: message})
}

func ResponseHTTPForbidden(w http.ResponseWriter, message string) {
	ResponseHTTP(w, http.StatusForbidden, ErrorHTTP{Message: message})
}

func ResponseHTTPNotFound(w http.ResponseWriter, message string) {
	ResponseHTTP(w, http.StatusNotFound, ErrorHTTP{Message: message})
}
