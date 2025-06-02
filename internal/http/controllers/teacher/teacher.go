package teacher

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/amito07/ems/internal/database"
	"github.com/amito07/ems/internal/models"
	"github.com/amito07/ems/internal/repository"
	"github.com/amito07/ems/internal/utils/response"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type TeacherController struct {
	teacherRepo *repository.TeacherRepository
	validator   *validator.Validate
}

func NewTeacherController() *TeacherController {
	db := database.GetDB()
	return &TeacherController{
		teacherRepo: repository.NewTeacherRepository(db),
		validator:   validator.New(),
	}
}

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Teacher Management System!"))
	}
}

func Create() http.HandlerFunc {
	controller := NewTeacherController()
	return func(w http.ResponseWriter, r *http.Request) {
		var teacher models.Teacher
		err := json.NewDecoder(r.Body).Decode(&teacher)

		if errors.Is(err, io.EOF) {
			response.WriteResponse(w, http.StatusBadRequest, "Empty body", response.GeneralErrorResponse(err))
			return
		}

		if err != nil {
			response.WriteResponse(w, http.StatusBadRequest, "Invalid request body", response.GeneralErrorResponse(err))
			return
		}

		// Validate the request body
		if err := controller.validator.Struct(teacher); err != nil {
			validateError := err.(validator.ValidationErrors)
			response.WriteResponse(w, http.StatusBadRequest, "Validation error", response.ValidationErrorResponse(validateError))
			return
		}

		// Check if email already exists
		existingTeacher, err := controller.teacherRepo.GetByEmail(teacher.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteResponse(w, http.StatusInternalServerError, "Database error", response.GeneralErrorResponse(err))
			return
		}
		if existingTeacher != nil {
			response.WriteResponse(w, http.StatusConflict, "Teacher with this email already exists", nil)
			return
		}

		// Create the teacher
		if err := controller.teacherRepo.Create(&teacher); err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, "Failed to create teacher", response.GeneralErrorResponse(err))
			return
		}

		response.WriteResponse(w, http.StatusCreated, "Teacher created successfully", teacher)
	}
}

func GetByID() http.HandlerFunc {
	controller := NewTeacherController()
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract ID from URL path
		path := strings.TrimPrefix(r.URL.Path, "/teachers/")
		idStr := strings.Split(path, "/")[0]
		
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.WriteResponse(w, http.StatusBadRequest, "Invalid teacher ID", response.GeneralErrorResponse(err))
			return
		}

		teacher, err := controller.teacherRepo.GetByID(uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.WriteResponse(w, http.StatusNotFound, "Teacher not found", nil)
				return
			}
			response.WriteResponse(w, http.StatusInternalServerError, "Database error", response.GeneralErrorResponse(err))
			return
		}

		response.WriteResponse(w, http.StatusOK, "Teacher retrieved successfully", teacher)
	}
}

func GetAll() http.HandlerFunc {
	controller := NewTeacherController()
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

		teachers, err := controller.teacherRepo.GetAll(offset, limit)
		if err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, "Failed to retrieve teachers", response.GeneralErrorResponse(err))
			return
		}

		// Get total count for pagination info
		total, err := controller.teacherRepo.Count()
		if err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, "Failed to count teachers", response.GeneralErrorResponse(err))
			return
		}

		result := map[string]interface{}{
			"teachers": teachers,
			"pagination": map[string]interface{}{
				"page":       page,
				"limit":      limit,
				"total":      total,
				"totalPages": (total + int64(limit) - 1) / int64(limit),
			},
		}

		response.WriteResponse(w, http.StatusOK, "Teachers retrieved successfully", result)
	}
}