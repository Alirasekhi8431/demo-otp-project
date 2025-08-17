package api

import (
	"net/http"

	"alirasekhi8431/demo-otp-project/internal/logic"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetOTP(c *gin.Context) {
	var requestBody struct {
		PhoneNumber string `json:"phoneNumber"`
		Username    string `json:"username"`
	}
	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "phone number is not correct"})
		return
	}
	str, err := logic.CreateOTPmsg(requestBody.PhoneNumber, requestBody.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "error processing"})
		logrus.Errorf("Error is => %v", err)
		return
	}
	logrus.Info(str)
	c.AbortWithStatusJSON(http.StatusAccepted, gin.H{"your otp is : ": str})

}

func CheckOtp(c *gin.Context) {
	var requestBody struct {
		Username string `json:"username"`
		Otp string `json:"otp"`
	}
	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "phone number is not correct"})
		return
	}
	logrus.Info("otp is => %v" , requestBody.Otp)
	ok , err := logic.CheckOtp(requestBody.Username , requestBody.Otp)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "otp is not correct" , "error" : err.Error()})
		return
	}
	  c.JSON(http.StatusOK, gin.H{"status": "ok"})


}

func SetupRoutes(router *gin.Engine) {
	router.POST("/getotp", GetOTP)
	router.POST("/check-otp" , CheckOtp)
}
