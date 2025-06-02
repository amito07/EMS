package repository

import (
	"github.com/amito07/ems/internal/models"
	"gorm.io/gorm"
)

type TeacherRepository struct {
	db *gorm.DB
}

func NewTeacherRepository(db *gorm.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

// Create creates a new teacher
func (r *TeacherRepository) Create(teacher *models.Teacher) error {
	return r.db.Create(teacher).Error
}

// GetByID retrieves a teacher by ID
func (r *TeacherRepository) GetByID(id uint) (*models.Teacher, error) {
	var teacher models.Teacher
	err := r.db.Preload("Courses").First(&teacher, id).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// GetByEmail retrieves a teacher by email
func (r *TeacherRepository) GetByEmail(email string) (*models.Teacher, error) {
	var teacher models.Teacher
	err := r.db.Where("email = ?", email).First(&teacher).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// GetByEmployeeID retrieves a teacher by employee_id
func (r *TeacherRepository) GetByEmployeeID(employeeID string) (*models.Teacher, error) {
	var teacher models.Teacher
	err := r.db.Where("employee_id = ?", employeeID).First(&teacher).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// GetAll retrieves all teachers with pagination
func (r *TeacherRepository) GetAll(offset, limit int) ([]models.Teacher, error) {
	var teachers []models.Teacher
	err := r.db.Offset(offset).Limit(limit).Preload("Courses").Find(&teachers).Error
	return teachers, err
}

// Update updates a teacher
func (r *TeacherRepository) Update(teacher *models.Teacher) error {
	return r.db.Save(teacher).Error
}

// Delete soft deletes a teacher
func (r *TeacherRepository) Delete(id uint) error {
	return r.db.Delete(&models.Teacher{}, id).Error
}

// Count returns the total number of teachers
func (r *TeacherRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Teacher{}).Count(&count).Error
	return count, err
}
