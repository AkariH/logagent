package main

import (
	"go.uber.org/zap"
	"logagent/config"
	"logagent/kafka"
	"logagent/utils/logger"
)

func test() {
	cfg := config.GetConfig()
	err := kafka.InitKafka([]string{cfg.KafkaConfig.Address})
	if err != nil {
		logger.Error("connect kafka failed", zap.Error(err))
	}
	logger.Infof("kafka client is %v", kafka.GetKafkaClient())

	logger.Info("test success")
	logger.Info("test success")
	logger.Info("test success")

}
