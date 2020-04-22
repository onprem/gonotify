package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Message sent successfully",
	})
}
