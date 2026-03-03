package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func Error(w http.ResponseWriter, status int, message string) {
	type jsonError struct {
		Error string `json:"error"`
	}

	erro := jsonError{
		Error: message,
	}

	buff, err := json.Marshal(erro)

	if err != nil {
		http.Error(w, "500: Internal server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(buff)))

	w.WriteHeader(status)

	w.Write(buff)
}

func JSON(w http.ResponseWriter, status int, data any) {

	buff, err := json.Marshal(data)

	if err != nil {
		Error(w, http.StatusInternalServerError, "500 Internal Server Error.")
		return
	}

	contentSize := strconv.Itoa(len(buff))

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", contentSize)
	w.WriteHeader(status)
	w.Write(buff)

}
