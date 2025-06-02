package rootrouter

import (
"net/http"

"github.com/amito07/ems/internal/http/routes/studentRoute"
"github.com/amito07/ems/internal/http/routes/teacherRoute"
)

func RouterInit() *http.ServeMux {
router := http.NewServeMux()
router.Handle("/health-check", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
res.Write([]byte("Server health is good!"))
}))
router.Handle("/api/v1/students/", studentroute.StudentMux())
router.Handle("/api/v1/teachers/", teacherroute.TeacherMux())
return router
}
