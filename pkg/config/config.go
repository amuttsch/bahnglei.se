package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type CountryConfig struct {
	Iso  string
	Name string
	Area string
}

type ThunderforestConfig struct {
	ApiKey   string
	Zoom     int
	MapStyle string
}

type Config struct {
	DatabaseUrl         string
	Countries           []CountryConfig
	ThunderforestConfig ThunderforestConfig
	OverpassUrl         string
}

func Read() *Config {
	conf := &Config{}

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.BindEnv("DatabaseUrl", "DATABASE_URL")
	viper.BindEnv("ThunderforestConfig.ApiKey", "THUNDERFOREST_API_KEY")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%v", err)
		panic("Invalid config")
	}

	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
		panic("Invalid config")
	}

	return conf
}
