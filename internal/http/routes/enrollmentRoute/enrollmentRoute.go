package enrollmentroute

import (
	"net/http"

	"github.com/amito07/ems/internal/http/controllers/enrollment"
)

func EnrollmentMux() http.Handler {
	enrollmentMux := http.NewServeMux()
	enrollmentMux.Handle("GET /list", enrollment.GetAllEnrollments())
	enrollmentMux.Handle("POST /create", enrollment.Create())
	enrollmentMux.Handle("GET /test", enrollment.Test())
	
	return http.StripPrefix("/api/v1/enrollments", enrollmentMux)
}
