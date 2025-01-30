package auth

import (
	"encoding/json"
	"net/http"

	apigateway "gorm_recruiter/api_gateway"
	"gorm_recruiter/handlers"
	"gorm_recruiter/models"
	"gorm_recruiter/seeders"

	"gorm.io/gorm"
)

func LogIn(env *models.Env, w http.ResponseWriter, r *http.Request) {
	//! Decode the incoming JSON body into a User struct
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response := models.Response{Message: err.Error(), Status: http.StatusBadRequest}
		handlers.SendResponse(w, response, http.StatusBadRequest)
		return
	}

	//! Encrypt the user's password
	encText, err := Encrypt(user.PasswordHash, seeders.PasswordHashingSecret)
	if err != nil {
		response := models.Response{Message: "Error encrypting password", Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}

	//! Check if the user exists in the database with the provided email and password
	var currentUser models.User
	if err := env.DB.Where("email = ? AND password_hash = ?", user.Email, encText).First(&currentUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := models.Response{Message: "Email ID or Password doesn't match", Status: http.StatusUnauthorized}
			handlers.SendResponse(w, response, http.StatusUnauthorized)
			return
		}
		response := models.Response{Message: err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}

	//! Generate a JWT token for the authenticated user
	tokenString, tokenError := apigateway.CreateToken(string(currentUser.UserID), currentUser.IsEmployer)
	if tokenError != nil {
		response := models.Response{Message: "Failed to create token", Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	response := models.Response{Message: "User exists", Status: http.StatusOK, Jwt: &tokenString}
	handlers.SendResponse(w, response, http.StatusOK)
}
