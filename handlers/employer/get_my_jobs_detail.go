package employer

import (
	"net/http"

	"gorm_recruiter/constants"
	"gorm_recruiter/handlers"
	"gorm_recruiter/models"

	"github.com/dgrijalva/jwt-go"
)

func GetMyJobsDetail(env *models.Env, w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(constants.Claims).(jwt.MapClaims)[constants.UniqueID].(string)
	type JobWithApplicants struct {
		models.Job
		Applicants []models.User `json:"applicants"`
	}
	var jobs []models.Job
	//! Fetch jobs posted by the user that are active
	if err := env.DB.Where("posted_by_id = ? AND is_active = ?", userID, true).Find(&jobs).Error; err != nil {
		response := models.Response{Message: "Error fetching job details!", Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	var jobsWithApplicants []JobWithApplicants
	//! Loop through each job and fetch the associated applicants
	for _, job := range jobs {
		var applicants []models.User
		if err := env.DB.Joins("JOIN job_applications ON job_applications.applicant_id = users.user_id").
			Where("job_applications.job_id = ?", job.JobID).
			Find(&applicants).Error; err != nil {
			response := models.Response{Message: "Error fetching applicants!", Status: http.StatusInternalServerError}
			handlers.SendResponse(w, response, http.StatusInternalServerError)
			return
		}
		//! Append job with applicants to the response slice
		jobsWithApplicants = append(jobsWithApplicants, JobWithApplicants{
			Job:        job,
			Applicants: applicants,
		})
	}
	response := models.Response{
		Message: "My Jobs fetched successfully!",
		Status:  http.StatusOK,
		Data:    jobsWithApplicants,
	}
	handlers.SendResponse(w, response, http.StatusOK)
}
