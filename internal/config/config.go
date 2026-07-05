package config

import (
	"time"

	"github.com/ezhigval/go-toolkit/config"
)

type Config struct {
	Port     string `env:"PORT" envDefault:"8081"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"json"`

	RedisAddr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB" envDefault:"0"`

	CacheTTL time.Duration `env:"CACHE_TTL" envDefault:"10m"`

	OpenWeatherAPIKey string        `env:"OPENWEATHER_API_KEY"`
	OpenWeatherURL    string        `env:"OPENWEATHER_URL" envDefault:"https://api.openweathermap.org/data/2.5/weather"`
	HTTPTimeout       time.Duration `env:"HTTP_TIMEOUT" envDefault:"5s"`

	BatchWorkers   int `env:"BATCH_WORKERS" envDefault:"4"`
	BatchQueueSize int `env:"BATCH_QUEUE_SIZE" envDefault:"32"`
}

func MustLoad() Config {
	return config.MustLoad[Config]()
}
