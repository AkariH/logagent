package main

import (
	"logagent/config"
	"logagent/kafka"
	"logagent/service"
	"logagent/tailfile"
	"logagent/utils/logger"
	"os"
	"os/signal"
	"syscall"
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
	producer, err := kafka.NewKafkaProducer([]string{cfg.KafkaConfig.Address})
	if err != nil {
		logger.Errorf("init kafka failed, err:%v", err)
	}
	logger.Infof("init kafka success")

	err = tailfile.Init(cfg.CollectConfig.LogPath)
	if err != nil {
		logger.Errorf("init tailfile failed, err:%v", err)
	}
	logger.Infof("init tailfile success")

	s := service.NewService(producer, service.KafkaConfiguration{
		Topic: cfg.KafkaConfig.Topic,
	})
    // Start the service
    if err := s.Start(); err != nil {
        logger.Errorf("Main: Service failed to start: %v", err)
        return
    }

    // Set up channel to listen for interrupt or terminate signals
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    // Block until a signal is received
    sig := <-sigs
    logger.Infof("Main: Received signal: %v. Initiating shutdown...", sig)

    // Initiate shutdown
    s.Shutdown()

    logger.Info("Main: Shutdown complete")
}
