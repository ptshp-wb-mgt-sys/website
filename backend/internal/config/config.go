// Package config contains the configuration for the application
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// server config
	Port string
	Env  string

	// URL for frontend
	FrontendURL string

	// Supabase config
	SupabaseURL        string
	SupabaseServiceKey string
	SupabaseJWTSecret  string
}

// LoadCfg loads the configuration from the environment
func LoadCfg() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	cfg := &Config{
		Port: getEnv("PORT", "3000"),
		Env:  getEnv("ENV", "development"),

		FrontendURL: getEnv("FRONTEND_URL", ""),

		SupabaseURL:        getEnv("SUPABASE_URL", ""),
		SupabaseServiceKey: getEnv("SUPABASE_SERVICE_KEY", ""),
		SupabaseJWTSecret:  getEnv("SUPABASE_JWT_SECRET", ""),
	}

	err := cfg.validateConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// getEnv grabs an environment variable or returns the default if not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

// validateConfig ensures all required environment variables are present
func (cfg *Config) validateConfig() error {
	var missingVars []string

	if cfg.SupabaseURL == "" {
		missingVars = append(missingVars, "SUPABASE_URL")
	}
	if cfg.SupabaseServiceKey == "" {
		missingVars = append(missingVars, "SUPABASE_SERVICE_KEY")
	}
	if cfg.SupabaseJWTSecret == "" {
		missingVars = append(missingVars, "SUPABASE_JWT_SECRET")
	}

	if len(missingVars) > 0 {
		return fmt.Errorf(
			"missing required environment variables: %s",
			strings.Join(missingVars, ", "),
		)
	}

	return nil
}
