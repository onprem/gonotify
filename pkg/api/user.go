package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/prmsrswt/gonotify/pkg/api/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/crypto/bcrypt"
)

func (api *API) queryUser(c *gin.Context) {
	l := log.With(api.logger, "route", "user")

	uID := int(c.MustGet("id").(float64))
	u := models.User{ID: uID}

	err := u.GetUserByID(api.DB)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid userID",
			})
			return
		}
		throwInternalError(c, l, err)
		return
	}

	c.JSON(http.StatusOK, u)
}

func (api *API) handleLogin(c *gin.Context) {
	l := log.With(api.logger, "route", "login")
	var user models.User

	type input struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var i input

	err := c.BindJSON(&i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "phone and password fields are required",
		})
		level.Error(l).Log("err", err)
		return
	}

	ph, err := parsePhone(i.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid phone number",
		})
		return
	}

	user.Phone = ph

	err = user.GetUserByPhone(api.DB)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid phone number or password",
			})
		} else {
			throwInternalError(c, l, err)
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
		throwInternalError(c, l, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   tokenStr,
	})
}

func (api *API) handleRegister(c *gin.Context) {
	l := log.With(api.logger, "route", "register")
	type input struct {
		Name     string `json:"name" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var i input

	var tempUser models.User
	err := c.BindJSON(&i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "all fields (name, phone, passowrd) are required",
		})
		return
	}
	level.Debug(l).Log("name", i.Name, "phone", i.Phone)

	tempUser.Name = i.Name
	tempUser.Verified = false

	ph, err := parsePhone(i.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid phone number",
		})
		return
	}

	i.Phone = ph
	wappNumber := fmt.Sprintf("whatsapp:%s", ph)
	level.Debug(l).Log("parsedPhone", ph, "wappNumber", wappNumber)

	tempUser.Phone = ph
	err = tempUser.GetUserByPhone(api.DB)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "phone number already exists",
		})
		return
	}
	if err != sql.ErrNoRows {
		throwInternalError(c, l, err)
		return
	}

	oldTime := time.Date(1950, time.January, 1, 0, 0, 0, 0, time.Local).Format(time.RFC3339)
	code := getVerificationCode()

	// People join the sandbox by sending a message, so we can risk it
	oldTime = time.Now().Format(time.RFC3339)

	hashByte, err := bcrypt.GenerateFromPassword([]byte(i.Password), bcrypt.DefaultCost)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}
	hash := string(hashByte)
	tempUser.Hash = hash

	level.Debug(l).Log("verification code", code)

	tx, err := api.DB.Begin()
	if err != nil {
		throwInternalError(c, l, err)
		return
	}
	defer tx.Rollback()

	uID, err := tempUser.InsertUser(tx)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}

	group := models.Group{Name: "default", UserID: int(uID)}
	gID, err := group.New(tx)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}

	num := models.Number{Phone: i.Phone, UserID: int(uID), Verified: false}

	nID, err := num.New(tx)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}

	_, err = tx.Exec(
		"INSERT INTO whatsappNodes(numberID, groupID, lastMsgReceived) VALUES(?, ?, ?)",
		nID,
		gID,
		oldTime,
	)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}

	nv := models.UserVerify{UserID: int(uID), NumberID: int(nID), Code: code}

	_, err = nv.New(tx)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}

	type TplInput struct {
		Code string
	}
	ti := TplInput{Code: code}

	var buf bytes.Buffer
	err = api.conf.VerifyTmpl.Execute(&buf, ti)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}
	level.Debug(l).Log("message", buf.String())

	err = api.TwilioClient.SendWhatsApp(api.conf.WhatsAppFrom, wappNumber, buf.String())
	if err != nil {
		throwInternalError(c, l, err)
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
	l := log.With(api.logger, "route", "userVerify")

	type input struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	var i input
	err := c.BindJSON(&i)
	if err != nil {
		throwInternalError(c, l, err)
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
	level.Debug(l).Log("phone", ph, "code", i.Code)

	user := &models.User{Phone: ph}
	err = user.GetUserByPhone(api.DB)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "cannot find user with given phone number",
			})
			return
		}
		throwInternalError(c, l, err)
		return
	}

	if user.Verified {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user is already verified",
		})
		return
	}

	verify := models.UserVerify{UserID: user.ID}

	err = verify.GetByUserID(api.DB)
	if err != nil {
		throwInternalError(c, l, err)
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
		throwInternalError(c, l, err)
		return
	}
	defer tx.Rollback()

	user.Verified = true
	_, err = user.UpdateUserByID(tx)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}

	n := models.Number{ID: verify.NumberID}
	err = n.GetByID(api.DB)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}
	n.Verified = true
	_, err = n.UpdateByID(tx)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}

	_, err = verify.DeleteByUserID(tx)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}

	tx.Commit()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 15).Unix(),
	})

	tokenStr, err := token.SignedString(api.conf.JWTSecret)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully verified your account",
		"token":   tokenStr,
	})
}
