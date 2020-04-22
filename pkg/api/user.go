package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// User models a user in database
type User struct {
	ID    int
	Name  string
	Email string
	Hash  string
}

func (api *API) handleLogin(c *gin.Context) {
	var user User

	type input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var i input

	err := c.BindJSON(&i)
	if err != nil {
		return
	}
	err = api.DB.QueryRow(
		"SELECT id, name, email, hash FROM users WHERE email = ?",
		i.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Hash)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid email or password",
			})
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	if i.Password != user.Hash {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 15).Unix(),
	})

	tokenStr, err := token.SignedString(api.JWTSecret)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   tokenStr,
	})
}

func (api *API) handleRegister(c *gin.Context) {
	type input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var i input

	var tempID int
	err := c.BindJSON(&i)
	if err != nil {
		return
	}

	err = api.DB.QueryRow("SELECT id FROM users WHERE email = ?", i.Email).Scan(&tempID)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email already exists",
		})
		return
	}
	if err != sql.ErrNoRows {
		c.Status(http.StatusInternalServerError)
		return
	}

	stmt, err := api.DB.Prepare("INSERT INTO users(name, email, hash) VALUES(?, ?, ?)")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	res, err := stmt.Exec(i.Name, i.Email, i.Password)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully registered",
		"user": gin.H{
			"id": id,
		},
	})
}
