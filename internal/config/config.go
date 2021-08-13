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
	PlanPath   string `mapstructure:"plan_path"`
	ReviewPath string `mapstructure:"review_path"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type UpdaterConfig struct {
	DownloadPath   string `mapstructure:"download_path"`
	SeleniumPath   string `mapstructure:"selenium_path"`
	BrowserMode    string `mapstructure:"browser_mode"`
	Port           int    `mapstructure:"port"`
	HuaweiLogin    string `mapstructure:"huawei_login"`
	HuaweiPassword string `mapstructure:"huawei_password"`
}

type Config struct {
	App     AppConfig     `mapstructure:"app"`
	Server  ServerConfig  `mapstructure:"server"`
	Redis   RedisConfig   `mapstructure:"redis"`
	Updater UpdaterConfig `mapstructure:"updater"`
}

func (conf *Config) PushToOSEnv() {
	os.Setenv("AUTH_SECRET", conf.Server.TokenSecret)
	os.Setenv("HUAWEI_LOGIN", conf.Updater.HuaweiLogin)
	os.Setenv("HUAWEI_PASS", conf.Updater.HuaweiPassword)
}
