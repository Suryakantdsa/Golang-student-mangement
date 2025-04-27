package storage

import "github/suryakantdsa/student-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents(limit int, skip int, queryParams map[string]string) (interface{}, error)
	UpdateStudent(id int64, body types.Student) (types.Student, error)
	DeleteStudent(id int64) (interface{}, error)
}
