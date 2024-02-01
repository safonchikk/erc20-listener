package util

import "github.com/spf13/viper"

type Config struct {
	EthNodeURL   string `mapstructure:"ETH_NODE_URL"`
	ContractAddr string `mapstructure:"CONTRACT_ADDR"`
	TokenDecimal string `mapstructure:"TOKEN_DECIMAL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
