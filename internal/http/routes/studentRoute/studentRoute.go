package studentroute

import (
	"net/http"

	"github.com/amito07/ems/internal/http/controllers/student"
)

func StudentMux() http.Handler{
	studentMux := http.NewServeMux()

	studentMux.Handle("GET /list", student.New())
	studentMux.Handle("POST /create", student.Create())
	return http.StripPrefix("/api/v1/students", studentMux)

}