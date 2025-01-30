package models

import (
	"net/http"

	"gorm.io/gorm"
)

type App struct {
	Router http.Handler
	Driver *gorm.DB
}

type Env struct {
	*gorm.DB
}
