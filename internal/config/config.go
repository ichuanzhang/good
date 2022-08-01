package config

import "github.com/spf13/viper"

var C = new(Settings)

func Setup(path string) (err error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err = v.ReadInConfig(); err != nil {
		return err
	}
	if err = v.Unmarshal(C); err != nil {
		return err
	}
	return checkConfig()
}

func checkConfig() error {

	return nil
}
