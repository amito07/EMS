package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/amito07/ems/internal/database"
	"github.com/amito07/ems/internal/models"
	"github.com/amito07/ems/internal/repository"
	"github.com/amito07/ems/internal/structure"
	passwordhashing "github.com/amito07/ems/internal/utils/passwordHashing"
	"github.com/amito07/ems/internal/utils/randomfunction"
	"github.com/amito07/ems/internal/utils/response"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type StudentController struct {
	studentRepo *repository.StudentRepository
	validator   *validator.Validate
}

func NewStudentController() *StudentController {
	db := database.GetDB()
	return &StudentController{
		studentRepo: repository.NewStudentRepository(db),
		validator:   validator.New(),
	}
}

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the EMS Student API!"))
	}
}

func Create() http.HandlerFunc {
	controller := NewStudentController()
	return func(w http.ResponseWriter, r *http.Request) {
		var student models.Student
		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteResponse(w, http.StatusBadRequest, "Empty body", response.GeneralErrorResponse(err))
			return
		}

		if err != nil {
			response.WriteResponse(w, http.StatusBadRequest, "Invalid request body", response.GeneralErrorResponse(err))
			return
		}

		// Validate the request body
		if err := controller.validator.Struct(student); err != nil {
			validateError := err.(validator.ValidationErrors)
			response.WriteResponse(w, http.StatusBadRequest, "Validation error", response.ValidationErrorResponse(validateError))
			return
		}

		// Generate student ID if not provided
		if student.S_ID == "" {
			// Simple ID generation - in production you'd want something more sophisticated
			student.S_ID = fmt.Sprintf("STU%03d", time.Now().Unix()%1000)
		}

		// Check if email already exists
		existingStudent, err := controller.studentRepo.GetByEmail(student.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteResponse(w, http.StatusInternalServerError, "Database error", response.GeneralErrorResponse(err))
			return
		}
		if existingStudent != nil {
			response.WriteResponse(w, http.StatusConflict, "Student with this email already exists", nil)
			return
		}

		// Create the student
		if err := controller.studentRepo.Create(&student); err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, "Failed to create student", response.GeneralErrorResponse(err))
			return
		}

		response.WriteResponse(w, http.StatusCreated, "Student created successfully", student)
	}
}

func GetByID() http.HandlerFunc {
	controller := NewStudentController()
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract ID from URL path
		path := strings.TrimPrefix(r.URL.Path, "/students/")
		idStr := strings.Split(path, "/")[0]
		
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.WriteResponse(w, http.StatusBadRequest, "Invalid student ID", response.GeneralErrorResponse(err))
			return
		}

		student, err := controller.studentRepo.GetByID(uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.WriteResponse(w, http.StatusNotFound, "Student not found", nil)
				return
			}
			response.WriteResponse(w, http.StatusInternalServerError, "Database error", response.GeneralErrorResponse(err))
			return
		}

		response.WriteResponse(w, http.StatusOK, "Student retrieved successfully", student)
	}
}

func GetAll() http.HandlerFunc {
	controller := NewStudentController()
	return func(w http.ResponseWriter, r *http.Request) {
		// Get pagination parameters
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

		page := 1
		limit := 10

		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		offset := (page - 1) * limit

		students, err := controller.studentRepo.GetAll(offset, limit)
		if err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, "Failed to retrieve students", response.GeneralErrorResponse(err))
			return
		}

		// Get total count for pagination info
		total, err := controller.studentRepo.Count()
		if err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, "Failed to count students", response.GeneralErrorResponse(err))
			return
		}

		result := map[string]interface{}{
			"students": students,
			"pagination": map[string]interface{}{
				"page":       page,
				"limit":      limit,
				"total":      total,
				"totalPages": (total + int64(limit) - 1) / int64(limit),
			},
		}

		response.WriteResponse(w, http.StatusOK, "Students retrieved successfully", result)
	}
}

func Update() http.HandlerFunc {
	controller := NewStudentController()
	 return func(w http.ResponseWriter, r *http.Request) {

		// Extract ID from URL path
		path := strings.TrimPrefix(r.URL.Path, "/students/update/")
		idStr := strings.Split(path, "/")[2]
		fmt.Println("ID to update:", idStr)

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.WriteResponse(w, http.StatusBadRequest, "Invalid student ID", response.GeneralErrorResponse(err))
			return
		}

		// Retrieve the student by ID
		isStudentExist, err := controller.studentRepo.GetByID(uint(id))

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.WriteResponse(w, http.StatusNotFound, "Student not found", nil)
				return
			}
			response.WriteResponse(w, http.StatusInternalServerError, "Database error", response.GeneralErrorResponse(err))
			return
		}

		fmt.Println("Existing student:", isStudentExist)

		var student models.Student
		decoderErr := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(decoderErr, io.EOF) {
			response.WriteResponse(w, http.StatusBadRequest, "Empty body", response.GeneralErrorResponse(err))
			return
		}

		if decoderErr != nil {
			response.WriteResponse(w, http.StatusBadRequest, "Invalid request body", response.GeneralErrorResponse(err))
			return
		}

		// Validate the request body
		if decoderErr := controller.validator.Struct(student); decoderErr != nil {
			validateError := err.(validator.ValidationErrors)
			response.WriteResponse(w, http.StatusBadRequest, "Validation error", response.ValidationErrorResponse(validateError))
			return
		}

		fmt.Println("Student to update:", student.Email)

		// cross check student email address
		if isStudentExist.Email != student.Email {
			response.WriteResponse(w, http.StatusBadRequest, "Email address not matched", nil)
			return
		}

		// update the student
		updateError := controller.studentRepo.UpdateMetaData(isStudentExist.ID, map[string]interface{}{"first_name": student.FirstName, "last_name": student.LastName})

		if updateError != nil {
			response.WriteResponse(w, http.StatusInternalServerError, "Failed to update student", response.GeneralErrorResponse(updateError))
			return
		}

		response.WriteResponse(w, http.StatusOK, "Student updated successfully",  nil)

	 }

}

func Delete() http.HandlerFunc {
	controller := NewStudentController()
	return func(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/students/delete/")
	idStr := strings.Split(path, "/")[2]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
	response.WriteResponse(w, http.StatusBadRequest, "Invalid student ID", response.GeneralErrorResponse(err))
	return
	}

	// Check if student exists
	_, err = controller.studentRepo.GetByID(uint(id))
	if err != nil {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.WriteResponse(w, http.StatusNotFound, "Student not found", nil)
		return
	}
	response.WriteResponse(w, http.StatusInternalServerError, "Database error", response.GeneralErrorResponse(err))
	return
	}

	// Delete the student
	if err := controller.studentRepo.Delete(uint(id)); err != nil {
	response.WriteResponse(w, http.StatusInternalServerError, "Failed to delete student", response.GeneralErrorResponse(err))
	return
	}

	response.WriteResponse(w, http.StatusOK, "Student deleted successfully", nil)
	}
}

// Business Logic for StudentController
func SignUp() http.HandlerFunc {
	controller := NewStudentController()
	return func(w http.ResponseWriter, r *http.Request) {
		var studentInfo structure.StudentSignupBody

		err := json.NewDecoder(r.Body).Decode(&studentInfo)

		if errors.Is(err, io.EOF) {
			response.WriteResponse(w, http.StatusBadRequest, "Empty body", response.GeneralErrorResponse(err))
			return
		}
		if err != nil {
			response.WriteResponse(w, http.StatusBadRequest, "Invalid request body", response.GeneralErrorResponse(err))
			return
		}

		fmt.Println("Student Info:", studentInfo)

		//check email already exists
		isStudentExist, err := controller.studentRepo.GetByEmail(studentInfo.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteResponse(w, http.StatusInternalServerError, "Database error", response.GeneralErrorResponse(err))
			return
		}
		if isStudentExist != nil {
			response.WriteResponse(w, http.StatusConflict, "Student with this email already exists", nil)
			return
		}

		fmt.Println("............. Okk")
		hashedPassword, err := passwordhashing.HashPassword(studentInfo.Password)
		if err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, "Failed to hash password", response.GeneralErrorResponse(err))
			return
		}
		fmt.Println("Hashed Password:", hashedPassword)

		// Create a new student model
		student := models.Student{
			FirstName: studentInfo.FirstName,
			LastName:  studentInfo.LastName,
			Email:     studentInfo.Email,
			S_ID:      randomfunction.GetStudentId(10),
			Password:  hashedPassword,
		}

		// Create the student in the database
		if err := controller.studentRepo.Create(&student); err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, "Failed to create student", response.GeneralErrorResponse(err))
			return
		}

		// Prepare the response
		response.WriteResponse(w, http.StatusCreated, "Student created successfully", nil)

	}
}

func Login() http.HandlerFunc {
	controller := NewStudentController()
	return func(w http.ResponseWriter, r *http.Request) {
		var loginInfo structure.StudentLoginBody

		err := json.NewDecoder(r.Body).Decode(&loginInfo)
		if errors.Is(err, io.EOF) {
			response.WriteResponse(w, http.StatusBadRequest, "Empty body", response.GeneralErrorResponse(err))
			return
		}

		if err != nil {
			response.WriteResponse(w, http.StatusBadRequest, "Invalid request body", response.GeneralErrorResponse(err))
			return
		}

		fmt.Println("Login Info:", loginInfo)

		// Retrieve the student by email
		student, err := controller.studentRepo.GetByEmail(loginInfo.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.WriteResponse(w, http.StatusNotFound, "Student not found", nil)
				return
			}
			response.WriteResponse(w, http.StatusInternalServerError, "Database error", response.GeneralErrorResponse(err))
			return
		}

		fmt.Println("Found Student:", student.Password)
		fmt.Println("Provided Password:", loginInfo.Password)
		// Verify the password
		if !passwordhashing.VerifyPassword(loginInfo.Password, student.Password) {
			response.WriteResponse(w, http.StatusUnauthorized, "Invalid credentials", nil)
			return
		}

		// generate jwt token
		token, err := randomfunction.GenerateJwtToken(student.Email, student.S_ID)
		if err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, "Failed to generate token", response.GeneralErrorResponse(err))
			return
		}

		responseData := structure.StudentLoginResponse{
			FirstName: student.FirstName,
			Token:     token,
		}

		response.WriteResponse(w, http.StatusOK, "Login successful", responseData)
	}
}