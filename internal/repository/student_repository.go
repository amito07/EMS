package repository

import (
	"github.com/amito07/ems/internal/models"
	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

// Create creates a new student
func (r *StudentRepository) Create(student *models.Student) error {
	return r.db.Create(student).Error
}

// GetByID retrieves a student by ID
func (r *StudentRepository) GetByID(id uint) (*models.Student, error) {
	var student models.Student
	err := r.db.Preload("Enrollments.Course").First(&student, id).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// GetByEmail retrieves a student by email
func (r *StudentRepository) GetByEmail(email string) (*models.Student, error) {
	var student models.Student
	err := r.db.Where("email = ?", email).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// GetByStudentID retrieves a student by student_id
func (r *StudentRepository) GetByStudentID(studentID string) (*models.Student, error) {
	var student models.Student
	err := r.db.Where("student_id = ?", studentID).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// GetAll retrieves all students with pagination
func (r *StudentRepository) GetAll(offset, limit int) ([]models.Student, error) {
	var students []models.Student
	err := r.db.Offset(offset).Limit(limit).Preload("Enrollments.Course").Find(&students).Error
	return students, err
}

// Update updates a student
func (r *StudentRepository) Update(student *models.Student) error {
	return r.db.Save(student).Error
}

// Delete soft deletes a student
func (r *StudentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Student{}, id).Error
}

// Count returns the total number of students
func (r *StudentRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Student{}).Count(&count).Error
	return count, err
}
