package main

import (
	dependencyInjection "parsing-service/dependency_injection"
	"time"

	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Kolkata")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	dependencyInjection.LoadDependecies()
}
