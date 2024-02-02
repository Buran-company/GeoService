package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type UserRepository interface {
	Create(db *sql.DB, sh Repo) error
	List(db *sql.DB, limit, offset int) ([]Repo, error)
	ConnectToDB() (*sql.DB, error)
}

type Repo struct {
	Query    string
	Request  string
	Response []byte
}

type DataBaseHandler struct {
	db *sql.DB
}

func (d *DataBaseHandler) ConnectToDB() error {
	err := godotenv.Load(filepath.Join("/home/hexedchild1/Kata/Repository/go-kata/course4Geoservice_1/", ".env"))
	if err != nil {
		return err
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	// connection string
	psqlconn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return err
	}

	// check db
	err = db.Ping()
	if err != nil {
		return err
	}

	d.db = db

	/*sqlStatement := `DROP TABLE geoservice`
	_, err = d.db.Exec(sqlStatement)
	if err != nil {
		return err
	}*/

	sqlStatement := `CREATE TABLE IF NOT EXISTS geoservice (
		id SERIAL PRIMARY KEY,
		Query text,
		Request text,
		Response bytea
	);`
	_, err = d.db.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return nil
}

func (d *DataBaseHandler) Create(sh Repo) error {
	sqlStatement := `
	INSERT INTO geoservice (Query, Request, Response)
	VALUES ($1, $2, $3);`
	_, err := d.db.Exec(sqlStatement, sh.Query, sh.Request, sh.Response)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBaseHandler) List(limit, offset int) ([]Repo, error) {
	var sqlStatement string
	var rows *sql.Rows
	var err error
	if limit == -1 && offset == -1 {
		sqlStatement = `
		SELECT * FROM geoservice;`
		rows, err = d.db.Query(sqlStatement)
	} else if limit != -1 && offset == -1 {
		sqlStatement = `
		SELECT * FROM geoservice LIMIT $1;`
		rows, err = d.db.Query(sqlStatement, limit)
	} else if limit == -1 && offset != -1 {
		sqlStatement = `
		SELECT * FROM geoservice OFFSET $1;`
		rows, err = d.db.Query(sqlStatement, offset)
	} else {
		sqlStatement = `
		SELECT * FROM geoservice LIMIT $1 OFFSET $2;`
		rows, err = d.db.Query(sqlStatement, limit, offset)
	}
	if err != nil {
		return []Repo{}, err
	}
	defer rows.Close()
	var shs []Repo
	for rows.Next() {
		var sh Repo
		var id = 0
		if err := rows.Scan(&id, &sh.Query, &sh.Request, &sh.Response); err != nil {
			return []Repo{}, err
		}
		shs = append(shs, sh)
	}
	if err := rows.Err(); err != nil {
		return []Repo{}, err
	}
	return shs, nil
}
