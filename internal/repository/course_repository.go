package repository

import (
	"github.com/amito07/ems/internal/models"
	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

// Create creates a new course
func (r *CourseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

// GetByID retrieves a course by ID
func (r *CourseRepository) GetByID(id uint) (*models.Course, error) {
	var course models.Course
	err := r.db.Preload("Teacher").Preload("Enrollments.Student").First(&course, id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// GetByCode retrieves a course by code
func (r *CourseRepository) GetByCode(code string) (*models.Course, error) {
	var course models.Course
	err := r.db.Where("course_code = ?", code).First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// GetByTeacherID retrieves courses by teacher ID
func (r *CourseRepository) GetByTeacherID(teacherID uint) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Where("teacher_id = ?", teacherID).Preload("Teacher").Find(&courses).Error
	return courses, err
}

// GetAll retrieves all courses with pagination
func (r *CourseRepository) GetAll(offset, limit int) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Offset(offset).Limit(limit).Preload("Teacher").Find(&courses).Error
	return courses, err
}

// Update updates a course
func (r *CourseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}

// Delete soft deletes a course
func (r *CourseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Course{}, id).Error
}

// Count returns the total number of courses
func (r *CourseRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Course{}).Count(&count).Error
	return count, err
}
