package postgress

import (
	"database/sql"
	"errors"
	"fmt"
	"github/suryakantdsa/student-api/internal/config"
	"github/suryakantdsa/student-api/internal/types"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS students(
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT,
			age INTEGER
		)`)
	if err != nil {
		return nil, err
	}
	return &Postgres{
		Db: db,
	}, nil

}

func (p *Postgres) CreateStudent(name string, email string, age int) (int64, error) {
	query := `INSERT INTO students (name, email, age) VALUES ($1, $2, $3) RETURNING id`

	var lastId int64
	err := p.Db.QueryRow(query, name, email, age).Scan(&lastId)
	if err != nil {
		return 0, err
	}

	return lastId, nil
}
func (p *Postgres) GetStudentById(id int64) (types.Student, error) {
	stmt, err := p.Db.Prepare("SELECT * FROM students where id = $1;")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student with id %s", fmt.Sprint(id))
		}

		return types.Student{}, fmt.Errorf("query error %w", err)
	}

	return student, nil
}

func (p *Postgres) GetStudents(limit int, skip int, params map[string]string) (interface{}, error) {
	if limit == 0 {
		limit = 20
	}
	if skip < 0 {
		skip = 0
	}

	query := "SELECT id,name,email,age FROM students"
	countQuery := "SELECT COUNT(*) FROM students"

	var args []interface{}
	var conditions []string

	if name, ok := params["name"]; ok && name != "" {
		conditions = append(conditions, "name ILIKE $"+strconv.Itoa(len(args)+1))
		args = append(args, "%"+name+"%")
	}
	if email, ok := params["email"]; ok && email != "" {
		conditions = append(conditions, "email = $"+strconv.Itoa(len(args)+1))
		args = append(args, email)
	}
	if age, ok := params["age"]; ok && age != "" {
		conditions = append(conditions, "age = $"+strconv.Itoa(len(args)+1))
		args = append(args, age)
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		query += whereClause
		countQuery += whereClause
	}

	if limit != -1 {
		query += " LIMIT $" + strconv.Itoa(len(args)+1)
		query += " OFFSET $" + strconv.Itoa(len(args)+2)
		args = append(args, limit, skip)
	}

	row, err := p.Db.Query(query, args...)
	if err != nil {
		return types.StudentListResponse{}, err
	}
	defer row.Close()

	students := []types.Student{}
	for row.Next() {
		var s types.Student
		if err := row.Scan(&s.Id, &s.Name, &s.Email, &s.Age); err != nil {
			return types.StudentListResponse{}, err
		}
		students = append(students, s)
	}

	if err := row.Err(); err != nil {
		return types.StudentListResponse{}, err
	}

	var total int
	if limit != -1 {
		err = p.Db.QueryRow(countQuery, args[:len(args)-2]...).Scan(&total)
	} else {
		err = p.Db.QueryRow(countQuery, args...).Scan(&total)
	}
	if err != nil {
		return types.StudentListResponse{}, err
	}

	if limit == -1 {
		return students, nil
	} else {
		return types.StudentListResponse{
			Total: total,
			Skip:  skip,
			Limit: limit,
			Data:  students,
		}, nil
	}
}
func (p *Postgres) UpdateStudent(id int64, body types.Student) (types.Student, error) {
	var setClaues []string
	var args []interface{}
	argsId := 1

	if body.Name != "" {
		setClaues = append(setClaues, "name= $"+strconv.Itoa(argsId))
		args = append(args, body.Name)
		argsId++
	}
	if body.Email != "" {
		setClaues = append(setClaues, "email= $"+strconv.Itoa(argsId))
		args = append(args, body.Email)
		argsId++
	}
	if body.Age != 0 {
		setClaues = append(setClaues, "age= $"+strconv.Itoa(argsId))
		args = append(args, body.Age)
		argsId++
	}

	if len(setClaues) == 0 {
		return types.Student{}, errors.New("no field to update")
	}

	query := `
	UPDATE students
	SET ` + strings.Join(setClaues, ", ") + `
	WHERE id = $` + strconv.Itoa(argsId) + `
	RETURNING id, name, email, age
`

	var updatedStudent types.Student
	args = append(args, id)
	err := p.Db.QueryRow(query, args...).Scan(&updatedStudent.Id, &updatedStudent.Name, &updatedStudent.Email, &updatedStudent.Age)

	if err != nil {
		return types.Student{}, err
	}

	return updatedStudent, nil
}
func (p *Postgres) DeleteStudent(id int64) (interface{}, error) {
	stmt, err := p.Db.Prepare("DELETE FROM students WHERE id = $1 RETURNING id, name, email, age;")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()
	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student with id %s", fmt.Sprint(id))
		}

		return types.Student{}, fmt.Errorf("query error %w", err)
	}

	return student, nil
}
