package apigateway

import (
	"fmt"
	"time"

	"gorm_recruiter/constants"
	"gorm_recruiter/seeders"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(uuid string, IsEmployer bool) (string, error) {
	//! Create a new token object
	var token *jwt.Token = jwt.New(jwt.SigningMethodHS256)
	//! Set claims (payload)
	var claims jwt.MapClaims = token.Claims.(jwt.MapClaims)
	claims[constants.UniqueID] = uuid                                 // Example data
	claims[constants.Expiry] = time.Now().Add(time.Hour * 700).Unix() // Token expires in 700 hours
	claims[constants.IsEmployer] = IsEmployer
	//! Generate encoded token and sign it with a secret
	tokenString, err := token.SignedString([]byte(seeders.JwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func VerifyToken(tokenString string, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//! Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}
