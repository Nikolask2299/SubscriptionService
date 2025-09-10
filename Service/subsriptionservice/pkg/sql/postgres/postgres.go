package postgres

import (
	"client"
	"database/sql"
	"fmt"
	"musicservice/interal/models"
	"strings"

	_ "github.com/lib/pq"
)

type Postgres struct {
    db *sql.DB
}

func NewPostgres(user, password, dbname, host, port string) (*Postgres, error) {
	psqlInfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)
    db, err := sql.Open("postgres", psqlInfo)
    if err!= nil {
        return nil, err
    }

    err = db.Ping()
    if err!= nil {
        return nil, err
    }

    return &Postgres{db: db}, nil
}


func (p *Postgres) Close() error {
    return p.db.Close()
}

