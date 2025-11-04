package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port                 string `envconfig:"PORT" required:"true"`
	MongoDBURI           string `envconfig:"MONGODB_URI" required:"true"`
	JWTSecret            string `envconfig:"JWT_SECRET" required:"true"`
	JWTExpiresIn         string `envconfig:"JWT_EXPIRES_IN" required:"true"`
	SMTPHost             string `envconfig:"SMTP_HOST" required:"true"`
	SMTPPort             int    `envconfig:"SMTP_PORT" required:"true"`
	EmailUser            string `envconfig:"EMAIL_USER" required:"true"`
	EmailPass            string `envconfig:"EMAIL_PASS" required:"true"`
	RedisURL             string `envconfig:"REDIS_URL" required:"true"`
	RedisProjectEditTTL  string `envconfig:"REDIS_PROJECT_EDIT_TTL" required:"true"`
	RedisProjectCacheTTL string `envconfig:"REDIS_PROJECT_CACHE_TTL" required:"true"`
	Timeout              string `envconfig:"TIMEOUT" required:"true"`
	CryptoAlgorithm      string `envconfig:"CRYPTO_ALGORITHM" required:"true"`
	KeyEncryptSecret     string `envconfig:"KEY_ENCRYPT_SECRET" required:"true"`
	IVLength             string `envconfig:"IV_LENGTH" required:"true"`
}

var Env Config

func Init() {

	_ = godotenv.Load("infra/.env")

	if err := envconfig.Process("", &Env); err != nil {
		log.Fatalf("env load error: %v", err)
	}
}
