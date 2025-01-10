package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"logagent/utils/logger"
)

type Config struct {
	KafkaConfig   KafkaConfig
	CollectConfig CollectConfig
}

type KafkaConfig struct {
	Address  string `yaml:"address"`
	Topic    string `yaml:"topic"`
	ChanSize int    `yaml:"chan_size"`
}

type CollectConfig struct {
	LogPath string `yaml:"log_collectPath"`
}

func InitConfig() {
	// Load config

	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal("error reading config file", zap.Error(err))
	}

}

func GetConfig() *Config {
	return &Config{
		KafkaConfig: KafkaConfig{
			Address: viper.GetString("kafka.address"),
			Topic:   viper.GetString("kafka.topic"),
		},
		CollectConfig: CollectConfig{
			LogPath: viper.GetString("logs.collectPath"),
		},
	}
}
