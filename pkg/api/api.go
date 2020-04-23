package api

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/prmsrswt/gonotify/pkg/twilio"
)

// API represents a API config object
type API struct {
	Host         string
	Port         string
	JWTSecret    []byte
	Gin          *gin.Engine
	WhatsAppFrom string
	TwilioClient *twilio.Twilio
	DB           *sql.DB
	logger       *log.Logger
}

// NewAPI creates a new API instance
func NewAPI(host, port, jwtSecret, twilioSID, twilioToken, whatsAppFrom string, db *sql.DB, logger *log.Logger) (*API, error) {
	err := bootstrapDB(db)
	if err != nil {
		return nil, err
	}

	return &API{
		Host:         host,
		Port:         port,
		JWTSecret:    []byte(jwtSecret),
		Gin:          gin.Default(),
		WhatsAppFrom: whatsAppFrom,
		TwilioClient: twilio.NewClient(twilioSID, twilioToken),
		DB:           db,
		logger:       logger,
	}, nil
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
		v1.GET("/ping", api.withAuth(), handlePing)

		v1.POST("/login", api.handleLogin)
		v1.POST("/register", api.handleRegister)
		v1.POST("/verify", api.handleUserVerify)

		v1.POST("/send", func(c *gin.Context) {
			c.Request.URL.Path = c.Request.URL.Path + "/whatsapp"
			api.Gin.HandleContext(c)
		})

		v1.POST("/send/whatsapp", api.withAuth(), api.handleWhatsApp)
	}
}

func (api *API) withAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("authorization")
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return api.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("email", claims["email"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
		}
	}
}

func handlePing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
		"email":   c.MustGet("email"),
	})
}
