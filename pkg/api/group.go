package api

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Group represents a group of nodes
type Group struct {
	ID            int            `json:"id"`
	Name          string         `json:"name"`
	UserID        int            `json:"userID"`
	WhatsAppNodes []WhatsAppNode `json:"whatsappNodes"`
}

func (api *API) queryGroups(c *gin.Context) {
	logger := log.With(api.logger, "route", "groups")

	uID := int(c.MustGet("id").(float64))

	wNodes := map[int][]WhatsAppNode{}

	wRows, err := api.DB.Query(
		`SELECT
			whatsappNodes.id,
			whatsappNodes.groupID,
			whatsappNodes.numberID,
			numbers.phone,
			whatsappNodes.lastMsgReceived
		FROM whatsappNodes
		LEFT JOIN numbers ON whatsappNodes.numberID = numbers.id
		WHERE numbers.userID = ?`,
		uID,
	)

	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}
	defer wRows.Close()

	for wRows.Next() {
		var n WhatsAppNode

		err = wRows.Scan(&n.ID, &n.GroupID, &n.NumberID, &n.Phone, &n.LastMsgReceived)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
			level.Error(logger).Log("err", err)
			return
		}

		wNodes[n.GroupID] = append(wNodes[n.GroupID], n)
	}

	groups := []Group{}

	rows, err := api.DB.Query(
		`SELECT id, userID, name FROM groups WHERE userID = ?`,
		uID,
	)

	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var g Group
		err = rows.Scan(&g.ID, &g.UserID, &g.Name)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
			level.Error(logger).Log("err", err)
			return
		}

		g.WhatsAppNodes = wNodes[g.ID]
		if g.WhatsAppNodes == nil {
			g.WhatsAppNodes = []WhatsAppNode{}
		}

		groups = append(groups, g)
	}

	c.JSON(http.StatusOK, gin.H{
		"groups": groups,
	})
}

func (api *API) handleAddGroup(c *gin.Context) {
	logger := log.With(api.logger, "route", "addGroup")

	type input struct {
		Name string `json:"name" binding:"required,alpha"`
	}

	var i input
	uID := int(c.MustGet("id").(float64))

	if err := c.ShouldBind(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "input validation failed",
		})
		return
	}

	i.Name = strings.ToLower(i.Name)

	var tmpID int
	err := api.DB.QueryRow(
		`SELECT id FROM groups WHERE name = ? AND userID = ?`,
		i.Name,
		uID,
	).Scan(&tmpID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "group already exists",
		})
		return
	}
	if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	_, err = api.DB.Exec(
		`INSERT INTO groups(userID, name) VALUES(?, ?)`,
		uID,
		i.Name,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
		level.Error(logger).Log("err", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "group created successfully",
	})
}
