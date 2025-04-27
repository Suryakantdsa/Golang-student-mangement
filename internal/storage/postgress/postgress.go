package postgress

import (
	"database/sql"
	"github/suryakantdsa/student-api/internal/config"

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
