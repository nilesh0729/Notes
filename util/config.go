package util

import "github.com/spf13/viper"

type config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	Secret        string `mapstructure:"PASSWORD"`
}

func LoadConfig(path string) (config config, err error) {

	viper.AddConfigPath(path)

	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err!=nil{
		return
	}
	err = viper.Unmarshal(&config)
	return
}
