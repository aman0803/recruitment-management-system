package employer

import (
	"net/http"

	"gorm_recruiter/handlers"
	"gorm_recruiter/models"

	"gorm.io/gorm"
)

func GetApplicantData(env *models.Env, w http.ResponseWriter, r *http.Request) {
	//! Get Applicant ID
	applicantID := r.URL.Query().Get("applicant_id")
	if applicantID == "" {
		response := models.Response{Message: "Applicant ID is required", Status: http.StatusBadRequest}
		handlers.SendResponse(w, response, http.StatusBadRequest)
		return
	}
	//! Fetch data from DB
	var applicant models.User
	if err := env.DB.Where("user_id = ?", applicantID).First(&applicant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := models.Response{Message: "Record Not found", Status: http.StatusNotFound}
			handlers.SendResponse(w, response, http.StatusNotFound)
			return
		}
		response := models.Response{Message: "Error fetching applicant", Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	response := models.Response{
		Message: "Applicant fetched successfully!",
		Status:  http.StatusOK,
		Data:    applicant,
	}
	handlers.SendResponse(w, response, http.StatusOK)
}
