package misc

import (
	"net/http"

	"gorm_recruiter/handlers"
	"gorm_recruiter/models"
)

func GetAllResumes(env *models.Env, w http.ResponseWriter, r *http.Request) {
	//! Fetch data from DB
	var resumes []models.Resume
	if err := env.DB.Preload("Educations").Preload("Experiences").Find(&resumes).Error; err != nil {
		response := models.Response{Message: "Error fetching resumes", Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	response := models.Response{
		Message: "resumes fetched successfully!",
		Status:  http.StatusOK,
		Data:    resumes,
	}
	handlers.SendResponse(w, response, http.StatusOK)
}
