package studentroute

import (
	"net/http"

	"github.com/amito07/ems/internal/http/controllers/student"
)

func StudentMux() http.Handler{
	studentMux := http.NewServeMux()

	studentMux.Handle("/list", student.New())
	studentMux.Handle("/create", student.Create())
	return http.StripPrefix("/api/v1/students", studentMux)

}