package teacherroute

import (
	"net/http"

	"github.com/amito07/ems/internal/http/controllers/teacher"
)

func TeacherMux() http.Handler {
	teacherMux := http.NewServeMux()

	teacherMux.Handle("GET /", teacher.GetAll())
	teacherMux.Handle("GET /list", teacher.New())
	teacherMux.Handle("POST /create", teacher.Create())
	teacherMux.Handle("GET /{id}", teacher.GetByID())
	
	return http.StripPrefix("/api/v1/teachers", teacherMux)
}
