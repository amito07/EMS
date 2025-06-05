-- Initialize EMS database schema
-- This script runs automatically when the PostgreSQL container starts for the first time

-- Create tables for EMS system
CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    student_id VARCHAR(50) UNIQUE NOT NULL,
    enrollment_date DATE DEFAULT CURRENT_DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS teachers (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    department VARCHAR(100),
    hire_date DATE DEFAULT CURRENT_DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS courses (
    id SERIAL PRIMARY KEY,
    course_code VARCHAR(20) UNIQUE NOT NULL,
    course_name VARCHAR(200) NOT NULL,
    credits INTEGER DEFAULT 3,
    teacher_id INTEGER REFERENCES teachers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS enrollments (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL REFERENCES students(id),
    course_id INTEGER NOT NULL REFERENCES courses(id),
    enrolled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active',
    grade VARCHAR(5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(student_id, course_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_students_email ON students(email);
CREATE INDEX IF NOT EXISTS idx_students_student_id ON students(student_id);
CREATE INDEX IF NOT EXISTS idx_teachers_email ON teachers(email);
CREATE INDEX IF NOT EXISTS idx_teachers_employee_id ON teachers(employee_id);
CREATE INDEX IF NOT EXISTS idx_courses_code ON courses(course_code);
CREATE INDEX IF NOT EXISTS idx_enrollments_student_course ON enrollments(student_id, course_id);


-- Insert sample data
INSERT INTO students (first_name, last_name, email, student_id) VALUES
    ('John', 'Doe', 'john.doe@email.com', 'STU001'),
    ('Jane', 'Smith', 'jane.smith@email.com', 'STU002'),
    ('Bob', 'Johnson', 'bob.johnson@email.com', 'STU003')
ON CONFLICT (email) DO NOTHING;

INSERT INTO teachers (first_name, last_name, email, employee_id, department) VALUES
    ('Dr. Sarah', 'Wilson', 'sarah.wilson@school.edu', 'TCH001', 'Computer Science'),
    ('Prof. Michael', 'Brown', 'michael.brown@school.edu', 'TCH002', 'Mathematics'),
    ('Dr. Lisa', 'Davis', 'lisa.davis@school.edu', 'TCH003', 'Physics')
ON CONFLICT (email) DO NOTHING;

INSERT INTO courses (course_code, course_name, credits, teacher_id) VALUES
    ('CS101', 'Introduction to Programming', 3, 1),
    ('MATH201', 'Calculus I', 4, 2),
    ('PHYS101', 'General Physics', 3, 3)
ON CONFLICT (course_code) DO NOTHING;

INSERT INTO enrollments (student_id, course_id, status) VALUES
    (1, 1, 'active'),
    (1, 2, 'active'),
    (2, 1, 'active'),
    (3, 3, 'completed')
ON CONFLICT (student_id, course_id) DO NOTHING;

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ems_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ems_user;
