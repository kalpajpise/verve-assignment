package config

import (
	"context"
	"log"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	ServiceName   string        `env:"SERVICE_NAME,default=verve-server"`
	JSONSizeLimit int64         `env:"JSON_SIZE_LIMIT,default=5242880"`
	HTTPHost      string        `env:"HTTP_HOST,default="`
	HTTPPort      string        `env:"HTTP_PORT,required"`
	HTTPTimeout   time.Duration `env:"HTTP_TIMEOUT,default=60s"`

	RedisEndpoint           string `env:"REDIS_ENDPOINT,default="`
	RedisPassword           string `env:"REDIS_PASSWORD,required"`
	RedisMaxRetries         int    `env:"REDIS_MAX_RETRIES,required"`
	RedisMinIdleConnections int    `env:"REDIS_MIN_IDLE_CONNECTIONS,required"`

	AwsSecretKey                string `env:"AWS_SECRET_KEY,required"`
	AwsAccessKey                string `env:"AWS_ACCESS_KEY,required"`
	AwsRegion                   string `env:"AWS_REGION,default=ap-south-1"`
	AwsKinesisUniqueCountStream string `env:"AWS_KINESIS_UNIQUE_COUNT_STREAM,required"`
}

func New() *Config {
	cfg := &Config{}

	if err := envconfig.Process(context.Background(), cfg); err != nil {
		log.Fatalf("Error while processing application config from env: %s", err)
	}

	return cfg
}
