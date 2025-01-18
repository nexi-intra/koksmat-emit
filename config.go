package main

import "github.com/spf13/viper"

func loadConfig() {
	// Load the configuration file
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()
}
