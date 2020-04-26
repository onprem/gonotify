package api

import (
	"database/sql"
	"net/http"
	"strings"
	"text/template"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/prmsrswt/gonotify/pkg/twilio"
)

// API represents a API config object
type API struct {
	conf         Config
	Gin          *gin.Engine
	TwilioClient *twilio.Twilio
	DB           *sql.DB
	logger       *log.Logger
}

// Config represents a configuration object for this API
type Config struct {
	TwilioSID             string
	TwilioToken           string
	Host                  string
	Port                  string
	JWTSecret             []byte
	WhatsAppFrom          string
	WebHookAccount        gin.Accounts
	VerifyTmpl, NotifTmpl *template.Template
}

// NewAPI creates a new API instance
func NewAPI(
	conf Config,
	db *sql.DB,
	logger *log.Logger,
) (*API, error) {

	err := bootstrapDB(db)
	if err != nil {
		return nil, err
	}

	return &API{
		conf:         conf,
		Gin:          gin.Default(),
		TwilioClient: twilio.NewClient(conf.TwilioSID, conf.TwilioToken),
		DB:           db,
		logger:       logger,
	}, nil
}

// Run is a wrapper around Register and Gin.Run()
func (api *API) Run() {
	api.Register()
	api.Gin.Run(strings.Join([]string{api.conf.Host, api.conf.Port}, ":"))
}

// Register creates all api endpoints in given instance of gin
func (api *API) Register() {
	v1 := api.Gin.Group("/api/v1")
	{
		v1.GET("/ping", api.withAuth(), handlePing)

		v1.GET("/user", api.withAuth(), api.queryUser)

		v1.POST("/login", api.handleLogin)
		v1.POST("/register", api.handleRegister)

		v1.POST("/verify", api.handleUserVerify)

		v1.POST("/send", api.withAuth(), func(c *gin.Context) {
			c.Request.URL.Path = c.Request.URL.Path + "/whatsapp"
			api.Gin.HandleContext(c)
		})

		v1.POST("/send/whatsapp", api.withAuth(), api.handleWhatsApp)

		v1.GET("/numbers", api.withAuth(), api.queryNumbers)
		v1.POST("/numbers/add", api.withAuth(), api.handleAddNumber)
		v1.POST("/numbers/verify", api.withAuth(), api.handleVerifyNumber)

		v1.GET("/groups", api.withAuth(), api.queryGroups)
		v1.POST("/groups/add", api.withAuth(), api.handleAddGroup)
		v1.POST("/groups/add/whatsapp", api.withAuth(), api.handleAddWhatsAppToGroup)

		v1.POST("/incoming", gin.BasicAuth(api.conf.WebHookAccount), api.handleIncoming)
	}
}

func (api *API) withAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("authorization")
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return api.conf.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("id", claims["id"])
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
		"id":      c.MustGet("id"),
	})
}
