package main

import "github.com/spf13/viper"

type AppConfig struct {
	Port string
	Ado  struct {
		PersonalAccessToken string
		Organization        string
		Project             string
		BaseUrl             string
	}
	Thunderdome struct {
		BaseUrl      string
		BattlePrefix string
	}
}

func getConfig() (*AppConfig, error) {
	config := AppConfig{}
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
