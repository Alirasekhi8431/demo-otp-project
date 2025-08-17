package models

import "time"

type Otp struct {
	Digits string
	TimeStamp time.Time
	Username string
	PhoneNumber string
}
