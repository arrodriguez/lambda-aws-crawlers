package main

import (
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("crawlers")
}

func main() {
	MakeRoutes()
}
