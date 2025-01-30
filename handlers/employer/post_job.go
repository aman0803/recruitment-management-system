package employer

import (
	"encoding/json"
	"net/http"
	"os/exec"

	"gorm_recruiter/constants"
	"gorm_recruiter/handlers"
	"gorm_recruiter/models"

	"github.com/dgrijalva/jwt-go"
)

func AddJob(env *models.Env, w http.ResponseWriter, r *http.Request) {
	//! Get user ID from jwt claims
	userID := r.Context().Value("claims").(jwt.MapClaims)[constants.UniqueID].(string)
	//! Decode the incoming JSON body into a Job struct
	var currentJob models.Job
	err := json.NewDecoder(r.Body).Decode(&currentJob)
	if err != nil {
		response := models.Response{Message: err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	//! Generating new id for job
	newUUID, err := exec.Command("uuidgen").Output()
	if err != nil {
		response := models.Response{Message: err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	currentJob.JobID = string(newUUID)
	currentJob.PostedByID = userID
	//! Add job in table
	if err := env.Create(&currentJob).Error; err != nil {
		response := models.Response{Message: err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	response := models.Response{Message: "Job Posted successfully!", Status: 200}
	handlers.SendResponse(w, response, http.StatusOK)
}
