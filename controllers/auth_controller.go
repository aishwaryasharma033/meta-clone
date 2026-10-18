package controllers

import (
	"net/http"

	"instagram-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	"time"
    "github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	DB *gorm.DB
}

type SignupRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (ctrl *AuthController) Signup(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Provide a valid username, email, and password of at least 8 characters",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not process password",
		})
		return
	}

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	if err := ctrl.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Username or email may already exist",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Signup successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	var user models.User

	if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
    return
}

	secret := os.Getenv("JWT_SECRET")

	 if secret == "" {
    c.JSON(http.StatusInternalServerError, gin.H{
        "error": "Server configuration error",
    })
    return

}

  if err := ctrl.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{
        "error": "Invalid email or password",
    })
    return
}

  if err := bcrypt.CompareHashAndPassword(
    []byte(user.PasswordHash),
    []byte(req.Password),
    ); err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{
        "error": "Invalid email or password",
    })
    return
}

claims := jwt.MapClaims{
    "user_id": user.ID,
    "exp":     time.Now().Add(24 * time.Hour).Unix(),
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

tokenString, err := token.SignedString([]byte(secret))
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "error": "Could not generate token",
    })
    return
}
c.SetCookie(
	"token",
	tokenString,
	86400,
	"/",
	"",
	false,
	true,
)
c.JSON(http.StatusOK, gin.H{
    "message":  "Login successful",
    "username": user.Username,
})

  
}