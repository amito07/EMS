package rootrouter

import (
	"net/http"

	enrollmentroute "github.com/amito07/ems/internal/http/routes/enrollmentRoute"
	studentroute "github.com/amito07/ems/internal/http/routes/studentRoute"
	teacherroute "github.com/amito07/ems/internal/http/routes/teacherRoute"
)

func RouterInit() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/health-check", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Server health is good!"))
	}))
	router.Handle("/api/v1/students/", studentroute.StudentMux())
	router.Handle("/api/v1/teachers/", teacherroute.TeacherMux())
	router.Handle("/api/v1/enrollments/", enrollmentroute.EnrollmentMux())
	return router
}
