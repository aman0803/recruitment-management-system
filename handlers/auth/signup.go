package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os/exec"

	apigateway "gorm_recruiter/api_gateway"
	"gorm_recruiter/handlers"
	"gorm_recruiter/models"
	"gorm_recruiter/seeders"
)

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

// ! This should be in an env file in production
func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// ! Encrypt method is to encrypt or hide any classified text
func Encrypt(text, secretKey string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

func SignUp(env *models.Env, w http.ResponseWriter, r *http.Request) {
	//! Decode incoming json body
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response := models.Response{Message: err.Error(), Status: http.StatusBadRequest}
		handlers.SendResponse(w, response, http.StatusBadRequest)
		return
	}
	//! Generate UUID
	newUUID, err := exec.Command("uuidgen").Output()
	if err != nil {
		response := models.Response{Message: err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	user.UserID = string(newUUID)
	//! Generate encrypted password
	encryptedPassword, err := Encrypt(user.PasswordHash, seeders.PasswordHashingSecret)
	if err != nil {
		response := models.Response{Message: err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	user.PasswordHash = encryptedPassword
	//! Create user
	if err := env.Create(&user).Error; err != nil {
		response := models.Response{Message: err.Error(), Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	//! Genrate token
	tokenString, tokenError := apigateway.CreateToken(string(newUUID), user.IsEmployer)
	if tokenError != nil {
		response := models.Response{Message: "Failed to create token", Status: http.StatusInternalServerError}
		handlers.SendResponse(w, response, http.StatusInternalServerError)
		return
	}
	response := models.Response{Message: "Created new user", Status: 200, Jwt: &tokenString}
	handlers.SendResponse(w, response, http.StatusOK)
}
