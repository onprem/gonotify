package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prmsrswt/gonotify/pkg/twilio"
)

type message struct {
	Group string `json:"group" binding:"-"`
	Body  string `json:"body" binding:"required"`
}

func (api *API) handleWhatsApp(c *gin.Context) {
	logger := log.With(*api.logger, "route", "whatsapp")
	var json message
	var groupID int

	uID := int(c.MustGet("id").(float64))

	if err := c.ShouldBind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if json.Group == "" {
		json.Group = "default"
	}

	err := api.DB.QueryRow(
		"SELECT id FROM groups WHERE userID = ? AND name = ?",
		uID,
		json.Group,
	).Scan(&groupID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}

	err = sendWhatsApp(api.DB, uID, groupID, json.Body, api.TwilioClient, api.WhatsAppFrom)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Message sent successfully",
	})
}

func sendWhatsApp(db *sql.DB, userID, groupID int, body string, tc *twilio.Twilio, from string) error {
	res, err := db.Exec(
		`INSERT INTO notifications(userID, groupID, body, timeSt) VALUES (?, ?, ?, ?)`,
		userID,
		groupID,
		body,
		time.Now().Format(time.RFC3339),
	)
	if err != nil {
		return err
	}

	notifID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	rows, err := db.Query(
		`SELECT id, phone, verified, lastMsgReceived FROM numbers WHERE groupID = ?`,
		groupID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	var errors []error
	type num struct {
		id    int
		phone string
	}
	var pendingNumbers []num

	for rows.Next() {
		var phone, lastMsgReceived string
		var id int
		var verified bool

		err := rows.Scan(&id, &phone, &verified, &lastMsgReceived)
		if err != nil {
			return err
		}

		phone = "whatsapp:" + phone

		if !verified {
			continue
		}

		last, err := time.Parse(time.RFC3339, lastMsgReceived)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		isAllowed := time.Since(last).Hours() < 24

		if isAllowed {
			err = tc.SendWhatsApp(from, phone, body)
			if err != nil {
				errors = append(errors, err)
				continue
			}
		} else {
			pendingNumbers = append(pendingNumbers, num{id, phone})
		}
	}

	for _, v := range pendingNumbers {
		_, err = db.Exec(
			`INSERT INTO pendingMsgs(notifID, numberID) VALUES(?, ?)`,
			notifID,
			v.id,
		)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		err = tc.SendWhatsApp(from, v.phone, `Your appointment is coming up on today at notification`) // Replace with actual Template
		if err != nil {
			errors = append(errors, err)
			continue
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("%d errors occured when sending messages", len(errors))
	}

	return nil
}
