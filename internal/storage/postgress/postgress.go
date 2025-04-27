package postgress

import (
	"database/sql"
	"fmt"
	"github/suryakantdsa/student-api/internal/config"
	"github/suryakantdsa/student-api/internal/types"
	"log/slog"

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

	stmt, err := p.Db.Prepare(`INSERT INTO students (name,email,age) VALUES ($1,$2,$3) RETURNING id`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, nil
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return lastId, nil

}
func (p *Postgres) GetStudentById(id int64) (types.Student, error) {
	slog.Info("id", id)
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
