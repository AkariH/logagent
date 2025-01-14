package kafka

import (
	"context"
	"fmt"
	"logagent/utils/logger"
	"sync"

	"github.com/IBM/sarama"
)

var (
	kafkaClient sarama.SyncProducer
	MsgChan     chan *sarama.ProducerMessage
	err         error
	producer    *KafkaProducer
)

type KafkaProducer struct {
	client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewKafkaProducer initializes and returns a new KafkaProducer
func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	client, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logger.Errorf("Failed to initialize Kafka producer: %v", err)
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	producer := &KafkaProducer{
		client:  client,
		msgChan: make(chan *sarama.ProducerMessage, 100), // Buffered channel
		ctx:     ctx,
		cancel:  cancel,
	}

	return producer, nil
}

func (p *KafkaProducer) Start() {
	logger.AddWithLog(&p.wg, 1, "KafkaProducer")
	go func() {
		defer logger.DoneWithLog(&p.wg, "KafkaProducer")
		for {
			select {
			case msg, ok := <-p.msgChan:
				if !ok {
					// Channel closed, terminate the goroutine
					logger.Debug("KafkaProducer: msgChan closed, terminating sender goroutine")
					return
				}
				partition, offset, err := p.client.SendMessage(msg)
				if err != nil {
					logger.Errorf("Failed to send message to Kafka: %v", err)
					// Optionally implement retry logic here
					continue
				}
				logger.Infof("Message sent to Kafka successfully. Partition: %d, Offset: %d, Message: %v", partition, offset, msg.Value)
			case <-p.ctx.Done():
				logger.Info("KafkaProducer: received shutdown signal")
				return
			}
		}
	}()
}

// Send enqueues a message to be sent to Kafka
func (p *KafkaProducer) Send(msg *sarama.ProducerMessage) error {
	select {
	case p.msgChan <- msg:
		logger.Debugf("Message enqueued to msgChan: %v", msg.Value)
		return nil
	default:
		logger.Errorf("msgChan is full. Failed to enqueue message: %v", msg.Value)
		return fmt.Errorf("message channel is full")
	}
}

// Close gracefully shuts down the KafkaProducer
func (p *KafkaProducer) Close() error {
	// Signal the sender goroutine to terminate
	p.cancel()
	// Close the msgChan to stop receiving new messages
	close(p.msgChan)
	// Wait for the sender goroutine to finish
	p.wg.Wait()
	// Close the Kafka client
	if err := p.client.Close(); err != nil {
		logger.Errorf("Failed to close Kafka producer client: %v", err)
		return err
	}
	logger.Info("KafkaProducer: closed successfully")
	return nil
}

func NewKafkaMsg(topic, msg string) *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}
}
