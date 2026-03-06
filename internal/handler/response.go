package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func JSON(w http.ResponseWriter, status int, data any) {
	buff, err := json.Marshal(data)

	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	contentSize := strconv.Itoa(len(buff))

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", contentSize)

	w.WriteHeader(status)

	w.Write(buff)
}
