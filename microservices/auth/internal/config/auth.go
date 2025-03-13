package config

import (
	"os"
	"time"
)

type AuthConfig struct {
	Salt       string
	SigningKey string

	TokenTTL time.Duration
}

func NewAuthConfig() *AuthConfig {
	return &AuthConfig{
		Salt: os.Getenv("salt"),
		SigningKey: os.Getenv("signingKey"),
		
		TokenTTL: 60 * time.Minute,
	}
}