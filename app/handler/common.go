package handler

import (
	"encoding/json"
	"net/http"
)

func respondSuccess(w http.ResponseWriter, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func respondError(w http.ResponseWriter, code int, message string) {
	payload := map[string]string{"error": message}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(payload)
	w.WriteHeader(code)
	w.Write([]byte(response))
}
