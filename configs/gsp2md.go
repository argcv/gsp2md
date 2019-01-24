package configs

import (
	"github.com/spf13/viper"
)

type Gsp2MdConfig struct {
	GoogleSpreadsheets struct {
		Client string `mapstructure:"client-id"`
		Secret string `mapstructure:"client-secret"`
	} `mapstructure:"gs"`
	Input []struct {
		Url    string   `mapstructure:"url"`
		Ranges []string `mapstructure:"ranges"`
	} `mapstructure:"in"`
	Output []struct {
		Type string `mapstructure:"type"`
		Name string `mapstructure:"name"`
	} `mapstructure:"out"`
}

func LoadGsp2MdConfig() (cfg *Gsp2MdConfig, err error) {
	cfg = &Gsp2MdConfig{}
	err = viper.Unmarshal(&cfg)
	return
}
