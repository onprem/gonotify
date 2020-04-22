package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prmsrswt/gonotify/pkg/twilio"
)

// API represents a API config object
type API struct {
	Host         string
	Port         string
	Gin          *gin.Engine
	WhatsAppFrom string
	TwilioClient *twilio.Twilio
}

// NewAPI creates a new API instance
func NewAPI(host, port, twilioSID, twilioToken, whatsAppFrom string) *API {
	return &API{
		Host:         host,
		Port:         port,
		Gin:          gin.Default(),
		WhatsAppFrom: whatsAppFrom,
		TwilioClient: twilio.NewClient(twilioSID, twilioToken),
	}
}

// Run is a wrapper around Register and Gin.Run()
func (api *API) Run() {
	api.Register()
	api.Gin.Run(strings.Join([]string{api.Host, api.Port}, ":"))
}

// Register creates all api endpoints in given instance of gin
func (api *API) Register() {
	v1 := api.Gin.Group("/api/v1")
	{
		v1.GET("/ping", handlePing)

		v1.POST("/send", func(c *gin.Context) {
			c.Request.URL.Path = c.Request.URL.Path + "/whatsapp"
			api.Gin.HandleContext(c)
		})

		v1.POST("/send/whatsapp", api.handleWhatsApp)
	}
}

func handlePing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
