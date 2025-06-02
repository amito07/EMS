package repository

import (
	"github.com/amito07/ems/internal/models"
	"gorm.io/gorm"
)

type EnrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

// Create creates a new enrollment
func (r *EnrollmentRepository) Create(enrollment *models.Enrollment) error {
	return r.db.Create(enrollment).Error
}

// GetByID retrieves an enrollment by ID
func (r *EnrollmentRepository) GetByID(id uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.db.Preload("Student").Preload("Course").First(&enrollment, id).Error
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

// GetByStudentAndCourse retrieves enrollment by student and course
func (r *EnrollmentRepository) GetByStudentAndCourse(studentID, courseID uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.db.Where("student_id = ? AND course_id = ?", studentID, courseID).First(&enrollment).Error
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

// GetByStudentID retrieves all enrollments for a student
func (r *EnrollmentRepository) GetByStudentID(studentID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	err := r.db.Where("student_id = ?", studentID).Preload("Course").Find(&enrollments).Error
	return enrollments, err
}

// GetByCourseID retrieves all enrollments for a course
func (r *EnrollmentRepository) GetByCourseID(courseID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	err := r.db.Where("course_id = ?", courseID).Preload("Student").Find(&enrollments).Error
	return enrollments, err
}

// GetAll retrieves all enrollments with pagination
func (r *EnrollmentRepository) GetAll(offset, limit int) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	err := r.db.Offset(offset).Limit(limit).Preload("Student").Preload("Course").Find(&enrollments).Error
	return enrollments, err
}

// Update updates an enrollment
func (r *EnrollmentRepository) Update(enrollment *models.Enrollment) error {
	return r.db.Save(enrollment).Error
}

// Delete soft deletes an enrollment
func (r *EnrollmentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Enrollment{}, id).Error
}

// Count returns the total number of enrollments
func (r *EnrollmentRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Enrollment{}).Count(&count).Error
	return count, err
}

// GetEnrollmentsByStatus retrieves enrollments by status
func (r *EnrollmentRepository) GetEnrollmentsByStatus(status string, offset, limit int) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	err := r.db.Where("status = ?", status).Offset(offset).Limit(limit).Preload("Student").Preload("Course").Find(&enrollments).Error
	return enrollments, err
}
