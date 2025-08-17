package logic

import (
	db "alirasekhi8431/demo-otp-project/internal/db"
	"alirasekhi8431/demo-otp-project/models"
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

func CheckOtp(username, otp string) (bool, error) {
	otps, err := db.GetUsersOTP(username)
	if err != nil {
		return false, err
	}
	if len(otps) > 3 {
		return false  , fmt.Errorf("Too many requests")
	}
	now := time.Now()
	twoMinOtps := make([]models.Otp , 0)
	cutoffTime := now.Add(-2 * time.Minute)
	for _, v := range otps {
		if !v.TimeStamp.Before(cutoffTime) {
			twoMinOtps = append(twoMinOtps, v)
		}
	}
	for _, v := range otps {
		if v.Digits == otp {
			return true, nil
		}
	}
	return false, fmt.Errorf("Error : the otp does not exist")

}
func CreateOTPmsg(phoneNumber, username string) (string, error) {
	user, err := db.GetUser(username)
	if err != nil {
		//Later fix this to return err
		return "", err
	}
	logrus.Info(user)
	//Check if the user is spamming
	otps, err := db.GetUsersOTP(username)
	if err != nil {
		return "", err
	}
	if len(otps) > 3 {
		return ""  , fmt.Errorf("Too many requests")
	}
	str, err := generateOTP(6)
	if err != nil {
		return "", err
	}
	otp := models.Otp{
		Username:    username,
		TimeStamp:   time.Now(),
		Digits:      str,
		PhoneNumber: phoneNumber,
	}
	err = db.InsertOTPForUser(otp)
	if err != nil {
		return "", err
	}
	return str, nil

}
func generateOTP(length int) (string, error) {
	const otpChars = "0123456789"
	otp := make([]byte, length)
	_, err := io.ReadAtLeast(rand.Reader, otp, length)
	if err != nil {
		return "", err
	}
	for i := 0; i < length; i++ {
		otp[i] = otpChars[int(otp[i])%len(otpChars)]
	}
	return string(otp), nil
}
