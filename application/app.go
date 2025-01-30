package application

import (
	"net/http"

	"gorm_recruiter/models"

	"gorm.io/gorm"
)

func New(driver *gorm.DB) *models.App {
	var d models.Env = models.Env{driver}
	var env *models.Env = &d
	return &models.App{Router: AppRoutes(env), Driver: driver}
}

func StartServer(app *models.App) error {
	server := &http.Server{Addr: ":8080", Handler: app.Router}
	err := http.ListenAndServe(server.Addr, server.Handler)
	/*
		or
		err := http.ListenAndServe(":8080", app.router)
	*/
	if err != nil {
		return err
	}
	return nil
}
