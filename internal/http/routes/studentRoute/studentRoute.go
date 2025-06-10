package studentroute

import (
	"net/http"

	"github.com/amito07/ems/internal/http/controllers/student"
)

func StudentMux() http.Handler {
	studentMux := http.NewServeMux()

	studentMux.Handle("GET /", student.GetAll())
	studentMux.Handle("GET /list", student.New())
	studentMux.Handle("POST /create", student.Create())
	studentMux.Handle("GET /{id}", student.GetByID())
	studentMux.Handle("PATCH /update/{id}", student.Update())
	studentMux.Handle("DELETE /delete/{id}", student.Delete())
	studentMux.Handle("POST /signup", student.SignUp())
	studentMux.Handle("POST /login", student.Login())
	
	return http.StripPrefix("/api/v1/students", studentMux)
}