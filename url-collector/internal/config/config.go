package config

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

// Config stores a config specific to this application
type Config struct {
	Nasa
	Port    uint `yaml:"port" env:"PORT" env-default:"8080"`
	Timeout uint `yaml:"timeout" env:"TIMEOUT" env-default:"30" validate:"required,min=5"`
}

// Nasa stores a config specific to nasa domain
type Nasa struct {
	URL                string `yaml:"url" env:"URL" env-default:"https://api.nasa.gov/planetary/apod"`
	ApiKey             string `yaml:"api_key" env:"API_KEY" env-default:"DEMO_KEY"`
	ConcurrentRequests uint   `yaml:"concurrent_requests" env:"CONCURRENT_REQUESTS" env-default:"5" validate:"required,min=1"`
}

// Load reads config file in yaml format and also validates the structure
func Load() (cfg Config, err error) {
	if err = cleanenv.ReadConfig("config.yaml", &cfg); err != nil {
		return
	}

	if err = validator.New().Struct(&cfg); err != nil {
		return
	}

	return
}

func (cfg Config) String() string {
	b, _ := json.Marshal(cfg)
	return string(b)
}
