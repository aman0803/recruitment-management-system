package main

import (
	"fmt"
	"log"

	"gorm_recruiter/application"
	"gorm_recruiter/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Started main function")
	dsn := "root:aman@0803@tcp(localhost:3306)/simplerest?charset=utf8&parseTime=True&loc=Local"
	driver, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
		return
	}
	fmt.Println("Database connected successfully:", driver)
	// Auto-migrate all models
	if err := driver.AutoMigrate(&models.User{}, &models.Job{}, &models.Resume{}, &models.Education{}, &models.Experience{}, &models.JobApplication{}); err != nil {
		log.Fatalf("failed to auto-migrate database: %v", err)
		return
	}
	var app *models.App = application.New(driver)
	if err := application.StartServer(app); err != nil {
		log.Fatal(err)
	}
}
