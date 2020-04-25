package api

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func (api *API) handleAddGroup(c *gin.Context) {
	logger := log.With(*api.logger, "route", "addGroup")

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
		c.Status(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}

	_, err = api.DB.Exec(
		`INSERT INTO groups(userID, name) VALUES(?, ?)`,
		uID,
		i.Name,
	)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		level.Error(logger).Log("err", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "group created successfully",
	})
}
