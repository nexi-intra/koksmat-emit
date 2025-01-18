package config

import "github.com/spf13/viper"

func Setup() {
	// Load the configuration file
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()
}
