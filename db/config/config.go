package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func init() {
	setDefault()
}

func setDefault() {
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.passwd", "abc123")
	viper.SetDefault("database.name", "test")
	viper.SetDefault("database.host", "localhost:3306")
	viper.SetDefault("database.connection", "tcp")
}

// GetInt returns a config property as int
func GetInt(property string) int {
	result := viper.GetInt(property)
	if result == 0 {
		fmt.Println("WARN: Property " + property + " is 0. Is this not set?")
	}
	return result
}

// GetString returns a config property as string
func GetString(property string) string {
	result := viper.GetString(property)
	if result == "" {
		fmt.Println("WARN: Property " + property + " is \"\". Is this not set?")
	}
	return result
}

// GetBool returns a config property as bool
func GetBool(property string) bool {
	return viper.GetBool(property)
}
