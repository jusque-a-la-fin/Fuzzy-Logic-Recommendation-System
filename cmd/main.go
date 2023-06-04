package main

import (
	"log"
	"vehicles/server"

	"vehicles/config"

	"github.com/spf13/viper"
)

func main() {

	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	if err := server.Run(viper.GetString("port")); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
