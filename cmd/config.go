package cmd

import (
	"github.com/nextrevision/traci/config"
	"github.com/spf13/viper"
	"log"
)

func getConfig() *config.Config {
	var c config.Config

	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return &c
}
