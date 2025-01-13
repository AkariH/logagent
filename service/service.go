package service

import (
	"github.com/IBM/sarama"
	"logagent/config"
	"logagent/kafka"
	"logagent/tailfile"
	"logagent/utils/logger"
)

var (
	msg *sarama.ProducerMessage
	cfg *config.Config
)

type Service struct {
    kafkaProducer *kafka.KafkaProducer
    KafkaConfig   KafkaConfiguration
}

type KafkaConfiguration struct {
    Topic string
    // Add other Kafka-related configurations if needed
}

// NewService creates a new Service instance
func NewService(producer *kafka.KafkaProducer, config KafkaConfiguration) *Service {
    return &Service{
        kafkaProducer: producer,
        KafkaConfig:   config,
    }
}

// ProcessLogs reads lines from the tailed file and sends them to Kafka
func (s *Service) ProcessLogs() error {
    for {
        line, err := tailfile.ReadLine()
        if err != nil {
            return err
        }

        logger.Infof("text: %s", line.Text)

        // Create a new Kafka message
        msg := kafka.NewKafkaMsg(s.KafkaConfig.Topic, line.Text)

        // Send the message using KafkaProducer
        s.kafkaProducer.Send(msg)
    }
}