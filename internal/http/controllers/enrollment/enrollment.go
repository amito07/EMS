package enrollment

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/amito07/ems/internal/database"
	"github.com/amito07/ems/internal/models"
	"github.com/amito07/ems/internal/repository"
	"github.com/amito07/ems/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

type EnrollmentController struct {
	enrollmentRepo *repository.EnrollmentRepository
	validator      *validator.Validate
}

func NewEnrollmentController() *EnrollmentController{
	db := database.GetDB()
	return &EnrollmentController{
		enrollmentRepo: repository.NewEnrollmentRepository(db),
		validator:      validator.New(),
	}
}

func Create() http.HandlerFunc {
	controller := NewEnrollmentController()
	return func (w http.ResponseWriter, r *http.Request) {
		fmt.Println(".............. HIT ...........................")
		var enrollment models.Enrollment
		err := json.NewDecoder(r.Body).Decode(&enrollment)

		if errors.Is(err, io.EOF) {
			response.WriteResponse(w, http.StatusBadRequest, "Empty body", response.GeneralErrorResponse(err))
			return
		}

		if err != nil {
			response.WriteResponse(w, http.StatusBadRequest, "Invalid request body", response.GeneralErrorResponse(err))
			return
		}

		// // Validate the request body
		// if err := controller.validator.Struct(enrollment); err != nil {
		// 	validateError := err.(validator.ValidationErrors)
		// 	response.WriteResponse(w, http.StatusBadRequest, "Validation error", response.ValidationErrorResponse(validateError))
		// 	return
		// }

		fmt.Println(enrollment)
		if err := controller.enrollmentRepo.Create(&enrollment); err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, "Failed to create enrollment", response.GeneralErrorResponse(err))
			return
		}

		response.WriteResponse(w, http.StatusCreated, "Enrollment created successfully", nil)
	}
}

func GetAllEnrollments() http.HandlerFunc {
	controller := NewEnrollmentController()
	return func(w http.ResponseWriter, r *http.Request){
		allEnrollments, err := controller.enrollmentRepo.GetAll(100)

		if err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, "Failed to fetch enrollments", response.GeneralErrorResponse(err))
			return
		}

		if len(allEnrollments) == 0 {
			response.WriteResponse(w, http.StatusNotFound, "No enrollments found", nil)
			return
		}

		  response.WriteResponse(w, http.StatusOK, "Enrollments fetched successfully", allEnrollments)
	}
 
}

func Test() http.HandlerFunc {
 controller := NewEnrollmentController()
 return func(w http.ResponseWriter, r *http.Request) {
  fmt.Println(".............. HIT ...........................")
  // Simulate some processing
  data, err := controller.enrollmentRepo.TestQuery()
  if err != nil {
	response.WriteResponse(w, http.StatusInternalServerError, "Failed to process request", response.GeneralErrorResponse(err))
	return
  }
  fmt.Println("Data from test query:", data)
  response.WriteResponse(w, http.StatusOK, "Test endpoint hit successfully", data)
 }
}