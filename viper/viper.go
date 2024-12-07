package viper

import (
	"github.com/spf13/viper"
)

type Config struct {
	C_User   string `mapstructure:"c_user"`
	Datr     string `mapstructure:"datr"`
	Fr       string `mapstructure:"fr"`
	Presence string `mapstructure:"presence"`
	Sb       string `mapstructure:"sb"`
	Wd       string `mapstructure:"wd"`
	Xs       string `mapstructure:"xs"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
