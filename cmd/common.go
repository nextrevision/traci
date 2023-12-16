package cmd

import (
	"fmt"
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

type EnumValue struct {
	Value   string
	Allowed []string
}

func (e *EnumValue) Set(val string) error {
	for _, allowedVal := range e.Allowed {
		if val == allowedVal {
			e.Value = val
			return nil
		}
	}
	return fmt.Errorf("invalid value '%s', allowed values are %v", val, e.Allowed)
}

func (e *EnumValue) Type() string {
	return "enum"
}

func (e *EnumValue) String() string {
	return e.Value
}
