package models

import (
	"time"
)

// Student represents a student in the EMS
type Student struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	FirstName      string    `json:"first_name" gorm:"column:first_name;not null" validate:"required"`
	LastName       string    `json:"last_name" gorm:"column:last_name;not null" validate:"required"`
	Email          string    `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	S_ID      string    `json:"s_id" gorm:"column:s_id;uniqueIndex;not null"`
	Password       *string    `json:"-" gorm:"column:password;not null" validate:"required,min=6"`
	EnrollmentDate *time.Time `json:"enrollment_date" gorm:"column:enrollment_date;type:date;default:CURRENT_DATE"`
	CreatedAt      time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relationships
	Enrollments []Enrollment `json:"enrollments,omitempty" gorm:"foreignKey:StudentID"`
}

// Teacher represents a teacher in the EMS
type Teacher struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	FirstName  string    `json:"first_name" gorm:"column:first_name;not null" validate:"required"`
	LastName   string    `json:"last_name" gorm:"column:last_name;not null" validate:"required"`
	Email      string    `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	EmployeeID string    `json:"employee_id" gorm:"column:employee_id;uniqueIndex;not null"`
	Department string    `json:"department" gorm:"size:100"`
	HireDate   *time.Time `json:"hire_date" gorm:"column:hire_date;type:date;default:CURRENT_DATE"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relationships
	Courses []Course `json:"courses,omitempty" gorm:"foreignKey:TeacherID"`
}

// Course represents a course in the EMS
type Course struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CourseCode string    `json:"course_code" gorm:"column:course_code;uniqueIndex;not null" validate:"required"`
	CourseName string    `json:"course_name" gorm:"column:course_name;not null" validate:"required"`
	Credits    int       `json:"credits" gorm:"default:3" validate:"required,min=1,max=10"`
	TeacherID  uint      `json:"teacher_id" gorm:"column:teacher_id"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relationships
	Teacher     Teacher      `json:"teacher,omitempty" gorm:"foreignKey:TeacherID;references:ID"`
	Enrollments []Enrollment `json:"enrollments,omitempty" gorm:"foreignKey:CourseID"`
}

// Enrollment represents student enrollment in courses
type Enrollment struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	StudentID      uint      `json:"student_id" gorm:"column:student_id;not null"`
	CourseID       uint      `json:"course_id" gorm:"column:course_id;not null"`
	EnrollmentDate *time.Time `json:"enrollment_date" gorm:"column:enrollment_date;type:date;default:CURRENT_DATE"`
	Grade          *string   `json:"grade" gorm:"size:5"`
	Status         string    `json:"status" gorm:"default:'active';size:20"` // active, completed, dropped
	CreatedAt      time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relationships
	Student Student `json:"student,omitempty" gorm:"foreignKey:StudentID;references:ID"`
	Course  Course  `json:"course,omitempty" gorm:"foreignKey:CourseID;references:ID"`
}

// TableName overrides the table name for custom naming
func (Student) TableName() string {
	return "students"
}

func (Teacher) TableName() string {
	return "teachers"
}

func (Course) TableName() string {
	return "courses"
}

func (Enrollment) TableName() string {
	return "enrollments"
}
