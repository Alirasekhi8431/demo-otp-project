package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	models "alirasekhi8431/demo-otp-project/models"
)

// To define usual CRUD on the database
var conn *pgx.Conn

func ConnectToDb(username, password, port string) {
	// urlExample := "postgres://username:password@localhost:5432/database_name"

	connStr := fmt.Sprintf("postgres://%v:%v@localhost:%v/my_db", username, password, port)
	var err error
	conn, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		logrus.Errorf("Error connecting to db %v", err)
	}
}

func InsertUser(username, password string) {
	checkQuery := `SELECT (username , password) FROM users WHERE username=$1`

	row := conn.QueryRow(context.Background(), checkQuery, username)
	var user models.User
	if err := row.Scan(&user.Username, user.Password); err != nil {
		if err == pgx.ErrNoRows {

		}else {
			logrus.Errorf("Username already exists.")
			return
		}
	}
	query := `INSERT INTO users (username , password) VALUES($1 , $2); `
	_, err := conn.Exec(context.Background(), query, username, password)
	if err != nil {
		logrus.Errorf("Error inserting => %v", err)
		return
	}
}
