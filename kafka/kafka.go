package kafka

import (
	"logagent/utils/logger"
	"sync"

	"github.com/IBM/sarama"
)

var (
	kafkaClient sarama.SyncProducer
	MsgChan     chan *sarama.ProducerMessage
	err         error
	producer	*KafkaProducer
)

type KafkaProducer struct {
    client  sarama.SyncProducer
    msgChan chan *sarama.ProducerMessage
    wg      sync.WaitGroup
}

func NewKafkaProducer(address []string) (*KafkaProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer(address, cfg)
	if err != nil {
        logger.Errorf("Failed to initialize Kafka producer: %v", err)
        return nil, err
    }

	logger.Infof("connect to kafka success")

    producer := &KafkaProducer{
        client:  client,
        msgChan: make(chan *sarama.ProducerMessage, 100), // Buffered channel
    }

    producer.wg.Add(1)
    go producer.start()

    return producer, nil
}

func GetKafkaClient() sarama.SyncProducer {
	return kafkaClient
}

func GetKafkaProducer() *KafkaProducer {
	return producer
}

// start listens to the msgChan and sends messages to Kafka
func (p *KafkaProducer) start() {
    defer p.wg.Done()
    for msg := range p.msgChan {
        partition, offset, err := p.client.SendMessage(msg)
        if err != nil {
            logger.Errorf("Failed to send message to Kafka: %v", err)
            // Optionally implement retry logic here
            continue
        }
        logger.Infof("Message sent to Kafka successfully. Partition: %d, Offset: %d, Message: %v", partition, offset, msg.Value)
    }
}

// Send enqueues a message to be sent to Kafka
func (p *KafkaProducer) Send(msg *sarama.ProducerMessage) {
    p.msgChan <- msg
}

// Close gracefully shuts down the KafkaProducer
func (p *KafkaProducer) Close() error {
    close(p.msgChan)
    p.wg.Wait()
    return p.client.Close()
}

func NewKafkaMsg(topic, msg string) *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}
}
