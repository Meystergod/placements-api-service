package config

import (
	"flag"
	"github.com/Meystergod/placements-api-service/internal/validator"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"strings"
	"sync"
)

// TODO: разобраться с конфигом, проверить и доработать его

type Config struct {
	AppConfig struct {
		IsDebug       bool   `yaml:"is-debug" env:"IS_DEBUG" env-default:"false"`
		IsDevelopment bool   `yaml:"is-development" env:"IS_DEV" env-default:"false"`
		LogLevel      string `yaml:"log-level" env:"LOG_LEVEL" env-default:"trace"`
	} `yaml:"app-config"`

	HTTP struct {
		IP       string   `yaml:"ip" env:"BIND_IP" env-default:"0.0.0.0"`
		Port     string   `env:"PORT" env-required:"true"`
		Partners []string `env:"PARTNERS" env-required:"true"`
		//ReadTimeout  time.Duration `yaml:"read-timeout" env:"HTTP-READ-TIMEOUT"`
		//WriteTimeout time.Duration `yaml:"write-timeout" env:"HTTP-WRITE-TIMEOUT"`
		CORS struct {
			AllowedMethods     []string `yaml:"allowed-methods" env:"HTTP-CORS-ALLOWED-METHODS"`
			AllowedOrigins     []string `yaml:"allowed-origins" env:"HTTP-CORS-ALLOWED-ORIGINS"`
			AllowCredentials   bool     `yaml:"allow-credentials" env:"HTTP-CORS-ALLOW-CREDENTIALS"`
			AllowedHeaders     []string `yaml:"allowed-headers" env:"HTTP-CORS-ALLOWED-HEADERS"`
			OptionsPassthrough bool     `yaml:"options-passthrough" env:"HTTP-CORS-OPTIONS-PASSTHROUGH"`
			ExposedHeaders     []string `yaml:"exposed-headers" env:"HTTP-CORS-EXPOSED-HEADERS"`
			Debug              bool     `yaml:"debug" env:"HTTP-CORS-DEBUG"`
		} `yaml:"cors"`
	} `yaml:"http"`
}

const pathToConfig = "configs/config.local.yaml"

var instance *Config
var once sync.Once

func GetArgs(cfg *Config) error {
	var port, partners string

	flag.StringVar(&port, "p", "", "port")
	flag.StringVar(&partners, "d", "", "partners")
	flag.Parse()

	partnersList := strings.Split(partners, ",")

	err := validator.ValidateArgs(port, partnersList)
	if err != nil {
		return err
	}

	cfg.HTTP.Port = port
	cfg.HTTP.Partners = partnersList

	return nil
}

func GetConfig() *Config {
	once.Do(func() {
		log.Print("initializing config")
		instance = &Config{}

		log.Print("getting application arguments for setting into config")
		if err := GetArgs(instance); err != nil {
			log.Fatal(err)
		}

		if err := cleanenv.ReadConfig(pathToConfig, instance); err != nil {
			helpDescription := "Help Config Description"
			help, _ := cleanenv.GetDescription(instance, &helpDescription)
			log.Print(help)
			log.Fatal(err)
		}
	})

	return instance
}
