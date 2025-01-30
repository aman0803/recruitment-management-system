package handlers

import (
	"encoding/json"
	"net/http"

	"gorm_recruiter/models"
)

func SendResponse(w http.ResponseWriter, response models.Response, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
