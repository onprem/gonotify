package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Notification represent a single message object
type Notification struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userID"`
	GroupID   int    `json:"groupID"`
	Body      string `json:"body"`
	TimeSt    string `json:"timeSt"`
	GroupName string `json:"groupName"`
}

func (api *API) queryNotifications(c *gin.Context) {
	logger := log.With(api.logger, "route", "notifications")

	uID := int(c.MustGet("id").(float64))

	notifications := []Notification{}

	rows, err := api.DB.Query(
		`SELECT notifications.id, notifications.userID, notifications.groupID, notifications.body,
		notifications.timeSt, groups.name FROM notifications
		LEFT JOIN groups ON notifications.groupID = groups.id
		WHERE notifications.userID = ?`,
		uID,
	)

	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var n Notification
		err = rows.Scan(&n.ID, &n.UserID, &n.GroupID, &n.Body, &n.TimeSt, &n.GroupName)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
			level.Error(logger).Log("err", err)
			return
		}

		notifications = append(notifications, n)
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
	})
}
