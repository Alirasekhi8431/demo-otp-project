package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	models "alirasekhi8431/demo-otp-project/models"
	"time"
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

func InsertUser(username, password string ) {
	checkQuery := `SELECT (username , password) FROM users WHERE username=$1`

	row := conn.QueryRow(context.Background(), checkQuery, username)
	var user models.User
	if err := row.Scan(&user); err != nil {
		if err == pgx.ErrNoRows {

		}else {
			logrus.Errorf("Username already exists. %v" , err)
			return
		}
	}
	query := `INSERT INTO users (username , password , registeration_date) VALUES($1 , $2,  $3); `
	_, err := conn.Exec(context.Background(), query, username, password , time.Now())
	if err != nil {
		logrus.Errorf("Error inserting => %v", err)
		return
	}
}
func GetUser(username string) (models.User , error) {
	checkQuery := `SELECT (username , password) FROM users WHERE username=$1`

	row := conn.QueryRow(context.Background(), checkQuery, username)
	var user models.User
	if err := row.Scan(&user); err != nil {
		
			return models.User{} , fmt.Errorf("User does not exist.")

	}
	return user , nil
}

func GetUsersOTP(username string) ([]models.Otp, error) {
	query := `
		SELECT otp, created_at FROM otps 
		WHERE username = $1 
		AND created_at >= NOW() - INTERVAL '10 minutes'
		ORDER BY created_at DESC;`

	rows, err := conn.Query(context.Background(), query, username)
	if err != nil {
		logrus.Errorf("Error getting user OTPs: %v", err)
		return nil, fmt.Errorf("error getting user OTPs")
	}
	defer rows.Close()

	var otps []models.Otp
	for rows.Next() {
		var otp models.Otp
		var createdAt time.Time
		if err := rows.Scan(&otp.Digits, &createdAt); err != nil {
			logrus.Errorf("Error scanning row: %v", err)
			continue
		}
		otp.Username = username
		otp.TimeStamp = createdAt
		otps = append(otps, otp)
	}

	if err := rows.Err(); err != nil {
		logrus.Errorf("Rows error: %v", err)
		return nil, fmt.Errorf("error with rows")
	}

	return otps, nil
}

func InsertOTPForUser(otp models.Otp) error{
	query := `INSERT INTO otps (username , otp ,created_at ,  phone_number ) VALUES($1 , $2, $3 ,  $4);`
	if _ , err := conn.Exec(context.Background() , query , otp.Username , otp.Digits , otp.TimeStamp , otp.PhoneNumber); err != nil {
		logrus.Errorf("Error inserting otp => %v" , err)
		return err
	}
	return nil


}