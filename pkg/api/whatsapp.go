package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log/level"
)

type message struct {
	Body string `json:"body" binding:"required"`
}

func (api *API) handleWhatsApp(c *gin.Context) {
	var json message
	if err := c.ShouldBind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := api.TwilioClient.SendWhatsApp(api.WhatsAppFrom, "whatsapp:+919950591608", json.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		level.Error(*api.logger).Log("err", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Message sent successfully",
	})
}
