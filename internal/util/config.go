package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	Secret              string        `mapstructure:"PASSWORD"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)

	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()
	
	// Explicitly bind environment variables for cloud deployments
	// This ensures Viper reads env vars even when .env file doesn't exist
	viper.BindEnv("DB_DRIVER")
	viper.BindEnv("DB_SOURCE")
	viper.BindEnv("SERVER_ADDRESS")
	viper.BindEnv("PASSWORD")
	viper.BindEnv("ACCESS_TOKEN_DURATION")

	// Try to read config file, but ignore if it doesn't exist
	// This allows the app to work with environment variables only (e.g., in cloud deployments)
	err = viper.ReadInConfig()
	if err != nil {
		// Ignore config file not found error, continue with environment variables
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			// Return only if it's a different error (not "file not found")
			return
		}
	}
	
	err = viper.Unmarshal(&config)
	return
}
