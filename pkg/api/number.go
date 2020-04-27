package api

import (
	"bytes"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Number models the numbers table in database
type Number struct {
	ID       int    `json:"id"`
	UserID   int    `json:"userID"`
	Phone    string `json:"phone"`
	Verified bool   `json:"verified"`
}

func (api *API) queryNumbers(c *gin.Context) {
	logger := log.With(api.logger, "route", "numbers")

	uID := int(c.MustGet("id").(float64))

	var numbers []Number

	rows, err := api.DB.Query(
		`SELECT id, userID, phone, verified FROM numbers WHERE userID = ?`,
		uID,
	)

	if err != nil && err != sql.ErrNoRows {
		c.Status(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var n Number
		err = rows.Scan(&n.ID, &n.UserID, &n.Phone, &n.Verified)

		if err != nil {
			c.Status(http.StatusInternalServerError)
			level.Error(logger).Log("err", err)
			return
		}

		numbers = append(numbers, n)
	}

	c.JSON(http.StatusOK, gin.H{
		"numbers": numbers,
	})
}

func (api *API) handleAddNumber(c *gin.Context) {
	logger := log.With(api.logger, "route", "addNumber")

	type input struct {
		Phone string `json:"phone" binding:"required"`
	}

	var i input
	uID := int(c.MustGet("id").(float64))

	if err := c.ShouldBind(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "all fields are required",
		})
		return
	}

	phone, err := parsePhone(i.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid phone number",
		})
		return
	}

	var tmpNum int
	err = api.DB.QueryRow(
		`SELECT id FROM numbers WHERE phone = ? AND userID = ?`,
		phone,
		uID,
	).Scan(&tmpNum)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "phone already registered",
		})
		return
	}
	if err != sql.ErrNoRows {
		c.Status(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}

	code := getVerificationCode()

	tx, err := api.DB.Begin()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}
	defer tx.Rollback()

	nRes, err := tx.Exec(
		`INSERT INTO numbers(userID, phone, verified) VALUES(?, ?, ?)`,
		uID,
		phone,
		0,
	)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}

	numID, err := nRes.LastInsertId()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}

	_, err = tx.Exec(
		`INSERT INTO numberVerify(numberID, code) VALUES(?, ?)`,
		numID,
		code,
	)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}

	type TplInput struct {
		Code string
	}
	ti := TplInput{Code: code}
	wappNumber := "whatsapp:" + phone

	var buf bytes.Buffer
	err = api.conf.VerifyTmpl.Execute(&buf, ti)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}
	level.Debug(logger).Log("message", buf.String())

	err = api.TwilioClient.SendWhatsApp(api.conf.WhatsAppFrom, wappNumber, buf.String())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "number added successfully",
	})
}

func (api *API) handleVerifyNumber(c *gin.Context) {
	logger := log.With(api.logger, "route", "verifyNumber")

	type input struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	var i input
	uID := int(c.MustGet("id").(float64))

	if err := c.ShouldBind(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "all fields are required",
		})
		return
	}

	phone, err := parsePhone(i.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid phone number",
		})
		return
	}

	var num Number
	err = api.DB.QueryRow(
		`SELECT id, verified FROM numbers WHERE phone = ? AND userID = ?`,
		phone,
		uID,
	).Scan(&num.ID, &num.Verified)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "given phone number doesn't exist",
			})
			return
		}
		c.Status(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}

	if num.Verified {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "phone number is already verified",
		})
		return
	}

	var verifyCode string
	err = api.DB.QueryRow(
		`SELECT code FROM numberVerify WHERE numberID = ?`,
		num.ID,
	).Scan(&verifyCode)

	if i.Code != verifyCode {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid verification code",
		})
		return
	}

	tx, err := api.DB.Begin()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`UPDATE numbers SET verified = ? WHERE id = ?`,
		1,
		num.ID,
	)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}

	_, err = tx.Exec(
		`DELETE FROM numberVerify WHERE numberID = ?`,
		num.ID,
	)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "phone number successfully verified",
	})
}
