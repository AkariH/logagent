package service

import (
	"context"
	"logagent/config"
	"logagent/kafka"
	"logagent/tailfile"
	"logagent/utils/logger"
	"sync"

	"github.com/IBM/sarama"
)

var (
	msg *sarama.ProducerMessage
	cfg *config.Config
)

type Service struct {
	kafkaProducer *kafka.KafkaProducer
	KafkaConfig   KafkaConfiguration
	wg            *sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
}

type KafkaConfiguration struct {
	Topic string
	// Add other Kafka-related configurations if needed
}

// NewService creates a new Service instance with initialized context and wait group
func NewService(producer *kafka.KafkaProducer, config KafkaConfiguration) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &Service{
		kafkaProducer: producer,
		KafkaConfig:   config,
		wg:            &sync.WaitGroup{},
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (s *Service) Start() error {
	// Start the KafkaProducer's internal goroutine
	s.kafkaProducer.Start()

	// Start the log processing in a separate goroutine
	logger.AddWithLog(s.wg, 1, "ProcessLogs")
	go func() {
		defer logger.DoneWithLog(s.wg, "ProcessLogs")
		if err := s.ProcessLogs(); err != nil {
			logger.Errorf("ProcessLogs encountered an error: %v", err)
			// Handle error as needed, e.g., initiate shutdown, retry, etc.
		}
	}()

	return nil
}

// ProcessLogs reads lines from the tailed file and sends them to Kafka
func (s *Service) ProcessLogs() error {
	for {
		select {
		case <-s.ctx.Done():
			logger.Info("Service: ProcessLogs: Context cancelled, stopping log processing")
			return nil
		default:
			line, err := tailfile.ReadLine()
			if err != nil {
				return err
			}
			if line == nil {
				logger.Debugf("line is nil or empty,skipping")
				continue
			}

			logger.Debugf("Service: Read line: %s", line.Text)

			// Create a new Kafka message
			msg := kafka.NewKafkaMsg(s.KafkaConfig.Topic, line.Text)

			logger.Debugf("Service: Created message: %+v", msg)

			// Send the message using KafkaProducer and handle potential errors
			if err := s.kafkaProducer.Send(msg); err != nil {
				logger.Errorf("Service: Failed to send message to Kafka: %v", err)
			}
		}
	}
}

// Shutdown gracefully shuts down the service, ensuring all goroutines finish execution
func (s *Service) Shutdown() {
	logger.Info("Service: Initiating shutdown...")

	// 1. Cancel the context to signal all goroutines to stop
	s.cancel()
	logger.Debug("Service: Context canceled")

	//// 3. Wait for all goroutines (e.g., ProcessLogs) to finish
	//s.wg.Wait()
	//logger.Debug("Service: All goroutines have finished")

	// 4. Close the KafkaProducer gracefully
	if err := s.kafkaProducer.Close(); err != nil {
		logger.Errorf("Service: Error closing KafkaProducer: %v", err)
	} else {
		logger.Info("Service: KafkaProducer closed successfully")
	}

	logger.Info("Service: Shutdown complete")
}
