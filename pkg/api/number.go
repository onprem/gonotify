package api

import (
	"bytes"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prmsrswt/gonotify/pkg/api/models"
)

func (api *API) queryNumbers(c *gin.Context) {
	l := log.With(api.logger, "route", "numbers")

	uID := int(c.MustGet("id").(float64))

	var numbers []models.Number

	rows, err := api.DB.Query(
		`SELECT numbers.id, numbers.userID, numbers.phone, numbers.verified, COUNT(whatsappNodes.groupID) as groups
		FROM numbers
		LEFT JOIN whatsappNodes ON numbers.id = whatsappNodes.numberID
		WHERE numbers.userID = ? GROUP BY numbers.id`,
		uID,
	)

	if err != nil && err != sql.ErrNoRows {
		throwInternalError(c, l, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var n models.Number
		err = rows.Scan(&n.ID, &n.UserID, &n.Phone, &n.Verified, &n.Groups)

		if err != nil {
			throwInternalError(c, l, err)
			return
		}

		numbers = append(numbers, n)
	}

	c.JSON(http.StatusOK, gin.H{
		"numbers": numbers,
	})
}

func (api *API) handleAddNumber(c *gin.Context) {
	l := log.With(api.logger, "route", "addNumber")

	type input struct {
		Phone string `json:"phone" binding:"required"`
	}

	var i input
	uID := int(c.MustGet("id").(float64))

	if err := c.ShouldBind(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "field 'phone' is required",
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

	num := models.Number{UserID: uID, Phone: phone}
	err = num.GetNumberByPhoneUserID(api.DB)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "phone already registered",
		})
		return
	}
	if err != sql.ErrNoRows {
		throwInternalError(c, l, err)
		return
	}

	code := getVerificationCode()

	tx, err := api.DB.Begin()
	if err != nil {
		throwInternalError(c, l, err)
		return
	}
	defer tx.Rollback()

	num.Verified = false
	numID, err := num.NewNumber(tx)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}

	_, err = tx.Exec(
		`INSERT INTO numberVerify(numberID, code) VALUES(?, ?)`,
		numID,
		code,
	)
	if err != nil {
		throwInternalError(c, l, err)
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

	c.JSON(http.StatusCreated, gin.H{
		"message": "number added successfully",
	})
}

func (api *API) handleVerifyNumber(c *gin.Context) {
	l := log.With(api.logger, "route", "verifyNumber")

	type input struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	var i input
	uID := int(c.MustGet("id").(float64))

	if err := c.ShouldBind(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "all fields (phone, code) are required",
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

	num := models.Number{UserID: uID, Phone: phone}
	err = num.GetNumberByPhoneUserID(api.DB)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "given phone number doesn't exist",
			})
			return
		}
		throwInternalError(c, l, err)
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
		throwInternalError(c, l, err)
		return
	}
	defer tx.Rollback()

	num.Verified = true
	_, err = num.UpdateNumberByID(tx)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}

	_, err = tx.Exec(
		`DELETE FROM numberVerify WHERE numberID = ?`,
		num.ID,
	)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "phone number successfully verified",
	})
}

func (api *API) handleRemoveNumber(c *gin.Context) {
	l := log.With(api.logger, "route", "removeNumber")

	type input struct {
		ID int `json:"id" binding:"required"`
	}

	var i input
	uID := int(c.MustGet("id").(float64))

	if err := c.ShouldBind(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "field 'id' is required",
		})
		return
	}

	num := models.Number{ID: i.ID}
	err := num.GetNumberByID(api.DB)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "given phone number doesn't exist",
			})
			return
		}
		throwInternalError(c, l, err)
		return
	}

	if num.UserID != uID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you are not authorized to delete this number",
		})
		return
	}

	user := models.User{ID: uID}
	err = user.GetUserByID(api.DB)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}

	if user.Phone == num.Phone {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot delete primary phone number",
		})
		return
	}

	tx, err := api.DB.Begin()
	if err != nil {
		throwInternalError(c, l, err)
		return
	}
	defer tx.Rollback()

	_, err = num.DeleteNumberByID(tx)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}

	_, err = tx.Exec(
		`DELETE FROM whatsappNodes WHERE numberID = ?`,
		i.ID,
	)
	if err != nil {
		throwInternalError(c, l, err)
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "phone number successfully deleted",
	})
}
