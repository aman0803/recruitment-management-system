package jobs

import (
	"net/http"

	"gorm_recruiter/handlers"
	"gorm_recruiter/models"

	"gorm.io/gorm"
)

func GetAllJobs(env *models.Env, w http.ResponseWriter, r *http.Request) {
	//! Fetch data from DB
	var jobs []models.Job
	if err := env.DB.Where("is_active = ?", true).Find(&jobs).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := models.Response{Message: "No Active Jobs!", Status: http.StatusNotFound}
			handlers.SendResponse(w, response, http.StatusNotFound)
			return
		}
		response := models.Response{Message: "Error fetching job details", Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	response := models.Response{
		Message: "Jobs fetched successfully!",
		Status:  http.StatusOK,
		Data:    jobs,
	}
	handlers.SendResponse(w, response, http.StatusOK)
}
