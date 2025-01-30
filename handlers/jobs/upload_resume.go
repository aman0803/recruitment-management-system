package jobs

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gorm_recruiter/constants"
	"gorm_recruiter/handlers"
	"gorm_recruiter/models"
	"gorm_recruiter/seeders"

	"github.com/dgrijalva/jwt-go"
)

func UploadResume(env *models.Env, w http.ResponseWriter, r *http.Request) {
	//! Get resume from request
	file, header, err := r.FormFile("resume")
	if err != nil {
		response := models.Response{Message: err.Error(), Status: http.StatusBadRequest}
		handlers.SendResponse(w, response, http.StatusBadRequest)
		return
	}
	//! save resume to directory
	resumeDir := "./resumes"
	if err := os.MkdirAll(resumeDir, os.ModePerm); err != nil {
		response := models.Response{Message: err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	// Create the file on the server
	resumeFilePath := filepath.Join(resumeDir, header.Filename)
	outFile, err := os.Create(resumeFilePath)
	if err != nil {
		response := models.Response{Message: err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	// Copy the uploaded file to the created file on the server
	if _, err := io.Copy(outFile, file); err != nil {
		response := models.Response{Message: "Error saving file: " + err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	// Read the saved file into a byte slice
	fileContent, err := os.ReadFile(resumeFilePath)
	if err != nil {
		response := models.Response{Message: "Error reading file: " + err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	//! Make call to external api for resume parsing
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.apilayer.com/resume_parser/upload", bytes.NewReader(fileContent))
	if err != nil {
		response := models.Response{Message: "Error creating request: " + err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("apikey", seeders.ResumeParserKey)
	response, err := client.Do(req)
	if err != nil {
		response := models.Response{Message: "Error uploading file: " + err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		response := models.Response{Message: "Error from resume parsing API: " + string(bodyBytes), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	//! Parse Resume
	var parsedResume models.Resume
	type ResumeResponse struct {
		Name       string              `json:"name"`
		Email      string              `json:"email"`
		Phone      string              `json:"phone"`
		Skills     []string            `json:"skills"`
		Education  []models.Education  `json:"education"`
		Experience []models.Experience `json:"experience"`
	}
	var parsedResponse ResumeResponse
	if err := json.NewDecoder(response.Body).Decode(&parsedResponse); err != nil {
		response := models.Response{Message: err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	newUUID, err := exec.Command("uuidgen").Output()
	if err != nil {
		response := models.Response{Message: err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	skillsString := strings.Join(parsedResponse.Skills, ", ")
	parsedResume.ResumeID = string(newUUID)
	parsedResume.UserID = r.Context().Value("claims").(jwt.MapClaims)[constants.UniqueID].(string)
	parsedResume.ResumeFileAddress = resumeFilePath
	parsedResume.Name = parsedResponse.Name
	parsedResume.Email = parsedResponse.Email
	parsedResume.Phone = parsedResponse.Phone
	parsedResume.Skills = skillsString
	parsedResume.Experiences = parsedResponse.Experience
	parsedResume.Educations = parsedResponse.Education
	if err := env.Create(&parsedResume).Error; err != nil {
		response := models.Response{Message: err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	message := models.Response{Message: "Resume Parsed successfully", Status: http.StatusOK, Data: parsedResume}
	handlers.SendResponse(w, message, http.StatusOK)
}
