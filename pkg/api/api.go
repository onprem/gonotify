package api

import "github.com/gin-gonic/gin"

// Register creates all api endpoints in given instance of gin
func Register(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.GET("/ping", handlePing)

		v1.POST("/send", func(c *gin.Context) {
			c.Request.URL.Path = "/v1/send/whatsapp"
			r.HandleContext(c)
		})

		v1.POST("/send/whatsapp", handlePing)
	}
}

func handlePing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
