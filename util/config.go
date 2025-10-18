package util

import (
	"github.com/spf13/viper"
)

// Config stores all the configuration settings for the application.
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	MigrationURL  string `mapstructure:"MIGRATION_URL"`
ServerAddress string `mapstructure:"SERVER_ADDRESS" default:"8081"`
}

// LoadConfig loads configuration settings from a file.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
