package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/crypto/bcrypt"
)

// User models a user in database
type User struct {
	ID       int
	Name     string
	Phone    string
	Hash     string
	Verified bool
}

// VerifyUser Models the verifyUser table
type VerifyUser struct {
	ID       int
	UserID   int
	NumberID int
	Code     string
}

func (api *API) queryUser(c *gin.Context) {
	logger := log.With(api.logger, "route", "user")

	var u User
	uID := int(c.MustGet("id").(float64))

	err := api.DB.QueryRow(
		`SELECT id, name, phone, verified FROM users WHERE id = ?`,
		uID,
	).Scan(&u.ID, &u.Name, &u.Phone, &u.Verified)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid userID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       u.ID,
		"name":     u.Name,
		"phone":    u.Phone,
		"verified": u.Verified,
	})
}

func (api *API) handleLogin(c *gin.Context) {
	logger := log.With(api.logger, "route", "login")
	var user User

	type input struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var i input

	err := c.BindJSON(&i)
	if err != nil {
		level.Error(logger).Log("err", err)
		return
	}

	ph, err := parsePhone(i.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid phone number",
		})
		return
	}

	i.Phone = ph

	err = api.DB.QueryRow(
		"SELECT id, name, phone, hash FROM users WHERE phone = ?",
		i.Phone,
	).Scan(&user.ID, &user.Name, &user.Phone, &user.Hash)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid phone number or password",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
			level.Error(logger).Log("err", err)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(i.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid phone number or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 15).Unix(),
	})

	tokenStr, err := token.SignedString(api.conf.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   tokenStr,
	})
}

func (api *API) handleRegister(c *gin.Context) {
	logger := log.With(api.logger, "route", "register")
	type input struct {
		Name     string `json:"name" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var i input

	var tempID int
	err := c.BindJSON(&i)
	if err != nil {
		return
	}
	level.Debug(logger).Log("name", i.Name, "phone", i.Phone)

	ph, err := parsePhone(i.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid phone number",
		})
		return
	}

	i.Phone = ph
	wappNumber := fmt.Sprintf("whatsapp:%s", ph)
	level.Debug(logger).Log("parsedPhone", ph, "wappNumber", wappNumber)

	err = api.DB.QueryRow("SELECT id FROM users WHERE phone = ?", i.Phone).Scan(&tempID)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "phone number already exists",
		})
		return
	}
	if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	oldTime := time.Date(1950, time.January, 1, 0, 0, 0, 0, time.Local).Format(time.RFC3339)
	code := getVerificationCode()

	hashByte, err := bcrypt.GenerateFromPassword([]byte(i.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}
	hash := string(hashByte)

	level.Debug(logger).Log("verification code", code)

	tx, err := api.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		"INSERT INTO users(name, phone, hash, verified) VALUES (?, ?, ?, ?)",
		i.Name,
		i.Phone,
		hash,
		0,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	uID, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	groupRes, err := tx.Exec("INSERT INTO groups(userID, name) VALUES(?, ?)", uID, "default")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	gID, err := groupRes.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	numRes, err := tx.Exec(
		"INSERT INTO numbers(phone, userID, verified) VALUES(?, ?, ?)",
		i.Phone,
		uID,
		0,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	nID, err := numRes.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	_, err = tx.Exec(
		"INSERT INTO whatsappNodes(numberID, groupID, lastMsgReceived) VALUES(?, ?, ?)",
		nID,
		gID,
		oldTime,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	_, err = tx.Exec(
		"INSERT INTO userVerify(userID, numberID, code) VALUES(?, ?, ?)",
		uID,
		nID,
		code,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	type TplInput struct {
		Code string
	}
	ti := TplInput{Code: code}

	var buf bytes.Buffer
	err = api.conf.VerifyTmpl.Execute(&buf, ti)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}
	level.Debug(logger).Log("message", buf.String())

	err = api.TwilioClient.SendWhatsApp(api.conf.WhatsAppFrom, wappNumber, buf.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully registered",
		"user": gin.H{
			"id": uID,
		},
	})
}

func (api *API) handleUserVerify(c *gin.Context) {
	logger := log.With(api.logger, "route", "userVerify")

	type input struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	var i input
	err := c.BindJSON(&i)
	if err != nil {
		level.Error(logger).Log("err", err)
		return
	}

	ph, err := parsePhone(i.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid phone number",
		})
		return
	}

	i.Phone = ph
	level.Debug(logger).Log("phone", ph, "code", i.Code)

	var tmpID int
	var tmpVerified bool
	err = api.DB.QueryRow(
		"SELECT id, verified FROM users WHERE phone = ?",
		i.Phone,
	).Scan(&tmpID, &tmpVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "cannot find user with given phone number",
			})
			return
		}
		level.Error(logger).Log("err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		return
	}

	if tmpVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user is already verified",
		})
		return
	}

	var verify VerifyUser
	err = api.DB.QueryRow(
		"SELECT id, userID, numberID, code FROM userVerify WHERE userID = ?",
		tmpID,
	).Scan(&verify.ID, &verify.UserID, &verify.NumberID, &verify.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	if i.Code != verify.Code {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid verification code",
		})
		return
	}

	tx, err := api.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE users SET verified = ? WHERE id = ?",
		1,
		tmpID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	_, err = tx.Exec(
		"UPDATE numbers SET verified = ? WHERE id = ?",
		1,
		verify.NumberID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	_, err = tx.Exec(
		"DELETE FROM userVerify WHERE userID = ?",
		tmpID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	tx.Commit()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":  tmpID,
		"exp": time.Now().Add(time.Hour * 24 * 15).Unix(),
	})

	tokenStr, err := token.SignedString(api.conf.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully verified your account",
		"token":   tokenStr,
	})
}
