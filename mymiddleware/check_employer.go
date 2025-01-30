package mymiddleware

import (
	"context"
	"fmt"
	"net/http"

	"gorm_recruiter/constants"
	"gorm_recruiter/handlers"
	"gorm_recruiter/models"

	"github.com/dgrijalva/jwt-go"
)

func CheckEmployer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claimsInterface := r.Context().Value(constants.Claims)
			if claimsInterface == nil {
				response := models.Response{Message: "Unauthorized request", Status: http.StatusUnauthorized}
				handlers.SendResponse(w, response, http.StatusUnauthorized)
				return
			}
			claims, ok := claimsInterface.(jwt.MapClaims)
			if !ok {
				response := models.Response{Message: "Invalid claims format", Status: http.StatusUnauthorized}
				handlers.SendResponse(w, response, http.StatusUnauthorized)
				return
			}
			if isEmployer, exists := claims[constants.IsEmployer].(bool); !exists || !isEmployer {
				response := models.Response{Message: "Not an employer!", Status: http.StatusUnauthorized}
				handlers.SendResponse(w, response, http.StatusUnauthorized)
				return
			} else {
				fmt.Println(claims[constants.IsEmployer].(bool))
				r = r.WithContext(context.WithValue(r.Context(), "claims", claims))
				next.ServeHTTP(w, r)
			}
		})
	}
}
