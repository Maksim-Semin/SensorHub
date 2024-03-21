package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type CardsData struct {
	Bytes   []byte `json:"bytes"`
	Command string `json:"command"`
	Message []byte `json:"message"`
}

var connStr = "host=localhost port=5432 user=postgres password=password dbname=CardsInfo sslmode=disable"

type Database interface {
	DBInit() error
	OpenConnection() (*sql.DB, error)
	CreateCard(uid []byte) error
	UpdateCard(uid []byte, command string) error
	DeleteCard(uid []byte) error
	GetCards() (CardsData, error)
}

type DB struct {
	db *sql.DB
}

func (d *DB) DBInit() error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}
	createTable := `
			CREATE TABLE IF NOT EXISTS Cards(
    		id SERIAL PRIMARY KEY,
    		bytes BYTEA NOT NULL UNIQUE,
    		command VARCHAR(255) NOT NULL,
			message BYTEA NOT NULL);
`
	_, err = db.Exec(createTable)
	if err != nil {
		return err
	}
	d.db = db
	return nil

}

func (d *DB) OpenConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return db, err
	}
	return db, nil
}

func (d *DB) CreateCard(uid []byte) error {
	db, err := d.OpenConnection()
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	_, err = db.Exec(`INSERT INTO Cards (bytes, command, message) VALUES ($1, $2, $3)`, uid, "UNDEFINED", "UNKNOWN")
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) UpdateCard(uid []byte, command string, message []byte) {
	db, err := d.OpenConnection()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE Cards SET command = $1, message = $2 WHERE bytes = $3", command, message, uid)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Change was success")
}

func (d *DB) GetCard(uid []byte) (CardsData, error) {
	var card CardsData

	db, err := d.OpenConnection()
	if err != nil {
		return CardsData{}, err
	}
	defer db.Close()

	row, err := db.Query("SELECT command, message FROM Cards WHERE bytes = $1", uid)
	if err != nil {
		return CardsData{}, err
	}
	for row.Next() {
		err := row.Scan(&card.Command, &card.Message)
		if err != nil {
			return CardsData{}, err
		}
	}
	return card, nil
}

func (d *DB) DeleteCard(uid []byte) error {
	db, err := d.OpenConnection()
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	_, err = db.Exec("DELETE from Cards WHERE bytes = $1", uid)
	if err != nil {
		return err
	}
	fmt.Println("Card was deleted")
	return nil
}

func (d *DB) GetCards() ([]CardsData, error) {
	var data []CardsData

	db, err := d.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT bytes, command  FROM Cards")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cardData CardsData
		err := rows.Scan(&cardData.Bytes, &cardData.Command)
		if err != nil {
			return nil, err
		}
		data = append(data, cardData)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}
