package config

import "os"

type ServerConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	TokenSecret    string `mapstructure:"token_secret"`
	PathToStatic   string `mapstructure:"static_path"`
	PathToHTMLGlob string `mapstructure:"html_glob_path"`
}

type AppConfig struct {
	DownloadPath string `mapstructure:"download_path"`
	PlanPath     string `mapstructure:"plan_path"`
	ReviewPath   string `mapstructure:"review_path"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Server ServerConfig `mapstructure:"server"`
	Redis  RedisConfig  `mapstructure:"redis"`
}

func (conf *Config) PushToOSEnv() {
	os.Setenv("AUTH_SECRET", conf.Server.TokenSecret)
}
