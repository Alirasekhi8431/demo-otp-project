package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"context"
)

func main() {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		gin.Logger()(c)
	})
	router.HTMLRender = ginview.Default()

	staticPath := fmt.Sprintf("%s/static", path.Base("."))
	router.Static("/static", staticPath)
	SetupRoutes(router)
	router.Run(":10000")
}



// Define a JWT secret key. In a real application, this should be a
// secure, random value stored as an environment variable.
var jwtKey = []byte("your_very_secret_key")

// CustomClaims defines the JWT token claims
type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func LogIn(c *gin.Context) {
	var requestBody struct {
		Username string `form:"username"`
		Otp      string `form:"otp"`
	}
	if err := c.ShouldBind(&requestBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "could not parse otp msg"})
		
		return
	}

	// Prepare and send the request to the /check-otp endpoint
	payload := gin.H{
		"username": requestBody.Username,
		"otp":      requestBody.Otp,
	}
	bodyBytes, _ := json.Marshal(payload)

	// Use the context to handle timeouts and cancellations for the HTTP request
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	reqUrl := "http://localhost:8080/check-otp"
	req, err := http.NewRequestWithContext(ctx, "POST", reqUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		logrus.Errorf("Failed to create request: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Failed to send request: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer resp.Body.Close()

	// Handle non-OK responses from the OTP check
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		logrus.Infof("OTP check failed with status: %s, body: %s", resp.Status, respBody)
		c.AbortWithStatusJSON(resp.StatusCode, gin.H{"msg": "OTP verification failed"})
		return
	}

	// OTP is verified, now generate a JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &CustomClaims{
		Username: requestBody.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		logrus.Errorf("Failed to generate JWT token: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Send the token back to the client
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"token":  tokenString,
	})
}

func ShowLogIn(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func SetupRoutes(router *gin.Engine){
	router.GET("/login" , ShowLogIn)
	router.POST("/login" , LogIn)
}
