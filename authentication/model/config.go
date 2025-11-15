package model

type Config struct {
	DatabaseURl   string `envconfig:"DATABASE_URL"`
	RedisAddr     string `envconfig:"REDIS_ADDR"`
	RedisPassword string `envconfig:"REDIS_PASSWORD"`
	SecretKey     string `envconfig:"SECRET_KEY"`
	FilePriPath   string `envconfig:"FILE_PRI_PATH"`
	FilePubPath   string `envconfig:"FILE_PUB_PATH"`
}
