package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/prmsrswt/gonotify/pkg/api"
)

type config struct {
	Server struct {
		Port string `yaml:"port" env:"PORT" env-default:"8080"`
		Host string `yaml:"host" env:"HOST" env-default:"0.0.0.0"`
	} `yaml:"server"`
	Twilio struct {
		SID          string `yaml:"sid" env:"TWILIO_SID"`
		Token        string `yaml:"token" env:"TWILIO_TOKEN"`
		WhatsAppFrom string `yaml:"whatsapp_from" env:"TWILIO_WHATSAPP_FROM"`
	} `yaml:"twilio"`
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

	gnAPI := api.NewAPI(
		cfg.Server.Host,
		cfg.Server.Port,
		cfg.Twilio.SID,
		cfg.Twilio.Token,
		cfg.Twilio.WhatsAppFrom,
	)
	gnAPI.Run()
}

func handleError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
