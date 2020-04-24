package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prmsrswt/gonotify/pkg/api"
)

type config struct {
	Server struct {
		Port      string `yaml:"port" env:"PORT" env-default:"8080"`
		Host      string `yaml:"host" env:"HOST" env-default:"0.0.0.0"`
		JWTSecret string `yaml:"jwt_secret" env:"JWT_SECRET"`
	} `yaml:"server"`
	Twilio struct {
		SID          string `yaml:"sid" env:"TWILIO_SID"`
		Token        string `yaml:"token" env:"TWILIO_TOKEN"`
		WhatsAppFrom string `yaml:"whatsapp_from" env:"TWILIO_WHATSAPP_FROM"`
		WebhookUser  string `yaml:"webhook_user" env:"TWILIO_WEBHOOK_USER"`
		WebhookPass  string `yaml:"webhook_password" env:"TWILIO_WEBHOOK_PASS"`
	} `yaml:"twilio"`
	Database struct {
		Path string `yaml:"path" env:"DATABASE_PATH" env-default:"gonotify.db"`
	} `yaml:"database"`
	Template struct {
		Verify       string `yaml:"verify" env:"TEMPLATE_VERIFY"`
		Notification string `yaml:"notification" env:"TEMPLATE_NOTIFICATION"`
	} `yaml:"templates"`
}

func main() {
	var cfg config
	var configPath string

	flag.StringVar(&configPath, "c", "config/config.yml", "path of config file")
	flag.Parse()

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		handleError(err)
	}

	db, err := sql.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		handleError(err)
	}

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)

	verifyTmpl := template.Must(template.New("verify").Parse(cfg.Template.Verify))
	notifTmpl := template.Must(template.New("notif").Parse(cfg.Template.Notification))

	conf := api.Config{
		TwilioSID:      cfg.Twilio.SID,
		TwilioToken:    cfg.Twilio.Token,
		Host:           cfg.Server.Host,
		Port:           cfg.Server.Port,
		JWTSecret:      []byte(cfg.Server.JWTSecret),
		WhatsAppFrom:   cfg.Twilio.WhatsAppFrom,
		WebHookAccount: gin.Accounts{cfg.Twilio.WebhookUser: cfg.Twilio.WebhookPass},
		VerifyTmpl:     verifyTmpl,
		NotifTmpl:      notifTmpl,
	}

	gnAPI, err := api.NewAPI(
		conf,
		db,
		&logger,
	)
	if err != nil {
		handleError(err)
	}
	gnAPI.Run()
}

func handleError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
