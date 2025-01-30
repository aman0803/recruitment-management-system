package application

import (
	"net/http"

	"gorm_recruiter/handlers"
	"gorm_recruiter/handlers/auth"
	"gorm_recruiter/handlers/employer"
	"gorm_recruiter/handlers/jobs"
	"gorm_recruiter/handlers/misc"
	"gorm_recruiter/models"
	"gorm_recruiter/mymiddleware"
	"gorm_recruiter/seeders"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func AppRoutes(env *models.Env) *chi.Mux {
	var router *chi.Mux = chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", DefaultRoute)
	router.Route("/", func(r chi.Router) {
		AuthRoutes(r, env)
	})
	router.Route("/employer", func(r chi.Router) {
		r.Use(mymiddleware.JwtVerify(seeders.JwtSecret))
		r.Use(mymiddleware.CheckEmployer())
		EmployerRoutes(r, env)
	})
	router.Route("/jobs", func(r chi.Router) {
		r.Use(mymiddleware.JwtVerify(seeders.JwtSecret))
		JobRoutes(r, env)
	})
	router.Route("/misc", func(r chi.Router) {
		r.Use(mymiddleware.JwtVerify(seeders.JwtSecret))
		GeneralRoutes(r, env)
	})
	return router
}

func DefaultRoute(w http.ResponseWriter, r *http.Request) {
	message := models.Response{Message: "Hey!", Status: 200}
	handlers.SendResponse(w, message, http.StatusOK)
}

func AuthRoutes(router chi.Router, env *models.Env) {
	router.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
		auth.SignUp(env, w, r)
	})
	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.LogIn(env, w, r)
	})
}

func EmployerRoutes(router chi.Router, env *models.Env) {
	router.Post("/postjob", func(w http.ResponseWriter, r *http.Request) {
		employer.AddJob(env, w, r)
	})
	router.Get("/getapplicantdata", func(w http.ResponseWriter, r *http.Request) {
		employer.GetApplicantData(env, w, r)
	})
	router.Get("/getmyjobsdetail", func(w http.ResponseWriter, r *http.Request) {
		employer.GetMyJobsDetail(env, w, r)
	})
}

func JobRoutes(router chi.Router, env *models.Env) {
	router.Get("/jobdata/{job_id}", func(w http.ResponseWriter, r *http.Request) {
		jobs.GetJobData(env, w, r)
	})
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		jobs.GetAllJobs(env, w, r)
	})
	router.Post("/uploadresume", func(w http.ResponseWriter, r *http.Request) {
		jobs.UploadResume(env, w, r)
	})
	router.Post("/apply/{job_id}", func(w http.ResponseWriter, r *http.Request) {
		jobs.ApplyToJob(env, w, r)
	})
}

func GeneralRoutes(router chi.Router, env *models.Env) {
	router.Get("/getall", func(w http.ResponseWriter, r *http.Request) {
		misc.GetAllApplicants(env, w, r)
	})
	router.Get("/getresumes", func(w http.ResponseWriter, r *http.Request) {
		misc.GetAllResumes(env, w, r)
	})
}
