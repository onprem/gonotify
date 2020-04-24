package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prmsrswt/gonotify/pkg/twilio"
)

func (api *API) handleWhatsApp(c *gin.Context) {
	logger := log.With(*api.logger, "route", "whatsapp")

	type message struct {
		Group string `json:"group" binding:"-"`
		Body  string `json:"body" binding:"required"`
	}

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

	err = sendWhatsApp(api.DB, uID, groupID, json.Body, api.TwilioClient, api.conf.WhatsAppFrom, api.conf.NotifTmpl)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Message sent successfully",
	})
}

func sendWhatsApp(db *sql.DB, userID, groupID int, body string, tc *twilio.Twilio, from string, tmpl *template.Template) error {
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

		type tmplData struct {
			Total int
		}
		var tData tmplData

		err = db.QueryRow(
			`SELECT COUNT(id) FROM pendingMsgs WHERE numberID = ?`,
			v.id,
		).Scan(&tData.Total)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		var msg bytes.Buffer
		err = tmpl.Execute(&msg, tData)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		err = tc.SendWhatsApp(from, v.phone, msg.String()) // Replace with actual Template
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

func (api *API) handleIncoming(c *gin.Context) {
	logger := log.With(*api.logger, "route", "incoming")

	type input struct {
		From string `form:"From"`
		Body string `form:"Body"`
	}
	var i input

	err := c.ShouldBind(&i)
	if err != nil {
		level.Error(logger).Log("err", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	number := strings.TrimPrefix(i.From, "whatsapp:")

	level.Debug(logger).Log("from", i.From, "body", i.Body)

	var numID int
	err = api.DB.QueryRow(
		`SELECT id FROM numbers WHERE phone = ?`,
		number,
	).Scan(&numID)
	if err != nil {
		if err == sql.ErrNoRows {
			level.Info(logger).Log("from", i.From, "msg", "messagse received from unknown number")
			c.Status(http.StatusNoContent)
			return
		}
		level.Error(logger).Log("err", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	_, err = api.DB.Exec(
		`UPDATE numbers SET lastMsgReceived = ? WHERE id = ?`,
		time.Now().Format(time.RFC3339),
		numID,
	)
	if err != nil {
		level.Error(logger).Log("err", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	type pending struct {
		id     int
		body   string
		timeSt string
	}

	rows, err := api.DB.Query(
		`SELECT pendingMsgs.id, notifications.body, notifications.timeSt FROM pendingMsgs
		INNER JOIN notifications ON pendingMsgs.notifID = notifications.id
		WHERE pendingMsgs.numberID = ?`,
		numID,
	)
	if err != nil {
		level.Error(logger).Log("err", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var msgs []pending
	body := "You have following notifications:"
	var ids []string

	for rows.Next() {
		var p pending
		err = rows.Scan(&p.id, &p.body, &p.timeSt)
		if err != nil {
			level.Error(logger).Log("err", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		msgs = append(msgs, p)
		body = body + "\n\n" + p.body
		ids = append(ids, strconv.Itoa(p.id))
	}

	if len(msgs) == 0 {
		body = "You have no new notifications"
	}

	err = api.TwilioClient.SendWhatsApp(api.conf.WhatsAppFrom, i.From, body)
	if err != nil {
		level.Error(logger).Log("err", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(msgs) > 0 {
		_, err := api.DB.Exec(
			`DELETE FROM pendingMsgs WHERE id IN (` + strings.Join(ids, ", ") + `)`,
		)
		if err != nil {
			level.Error(logger).Log("err", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.Status(http.StatusNoContent)
}
