package main

import (
	"logagent/config"
	"logagent/kafka"
	"logagent/service"
	"logagent/tailfile"
	"logagent/utils/logger"
)

//log collect client

func main() {

	// init
	// 1. read config
	// read logs via tail
	// send logs to kafka via sarama
	Init()
}

func Init() {

	config.InitConfig()
	cfg := config.GetConfig()
	defer logger.Sync()
	_,err := kafka.NewKafkaProducer([]string{cfg.KafkaConfig.Address})
	if err != nil {
		logger.Errorf("init kafka failed, err:%v", err)
	}
	logger.Infof("init kafka success")
	err = tailfile.Init(cfg.CollectConfig.LogPath)
	if err != nil {
		logger.Errorf("init tailfile failed, err:%v", err)
	}
	logger.Infof("init tailfile success")

	s:= service.NewService(kafka.GetKafkaProducer(), service.KafkaConfiguration{
		Topic: cfg.KafkaConfig.Topic,
})
s.ProcessLogs()

}