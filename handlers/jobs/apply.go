package jobs

import (
	"net/http"
	"os/exec"

	"gorm_recruiter/constants"
	"gorm_recruiter/handlers"
	"gorm_recruiter/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func ApplyToJob(env *models.Env, w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, constants.JobID)
	userID := r.Context().Value("claims").(jwt.MapClaims)[constants.UniqueID].(string)
	//! Check if job exists
	var job models.Job
	if err := env.DB.Where("job_id = ?", jobID).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := models.Response{Message: "Job doesn't exist or has been deleted!", Status: http.StatusNotFound}
			handlers.SendResponse(w, response, http.StatusNotFound)
			return
		}
		response := models.Response{Message: "Error checking job existence", Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	//! Check if user has already applied
	var existingApplication models.JobApplication
	if err := env.DB.Where("applicant_id = ? AND job_id = ?", userID, jobID).First(&existingApplication).Error; err == nil {
		response := models.Response{Message: "You have already applied for this job!", Status: http.StatusConflict}
		handlers.SendResponse(w, response, http.StatusConflict)
		return
	} else if err != gorm.ErrRecordNotFound {
		response := models.Response{Message: "Error checking application status", Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	//! Generate Application ID
	applicationID, err := exec.Command("uuidgen").Output()
	if err != nil {
		response := models.Response{Message: "Error generating Application ID", Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	var application models.JobApplication
	application.ApplicationID = string(applicationID)
	application.ApplicantID = userID
	application.JobID = jobID
	//! Add application to table
	if err := env.Create(&application).Error; err != nil {
		response := models.Response{Message: err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	response := models.Response{Message: "Job Applied successfully!", Status: 200}
	handlers.SendResponse(w, response, http.StatusOK)
}
