package kafka

import (
	"github.com/IBM/sarama"
	"logagent/utils/logger"
)

var (
	kafkaClient sarama.SyncProducer
	MsgChan     chan *sarama.ProducerMessage
	err         error
)

func InitKafka(address []string, chanSize int) error {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.Return.Successes = true

	kafkaClient, err = sarama.NewSyncProducer(address, cfg)
	if err != nil {
		logger.Errorf("connect to kafka failed, err:%v", err)
	}

	logger.Infof("connect to kafka success")

	MsgChan = make(chan *sarama.ProducerMessage, chanSize)
	go sendMsgToKafka()
	return nil
}

func GetKafkaClient() sarama.SyncProducer {
	return kafkaClient
}

// read msg from MsgChan, send to kafka
func sendMsgToKafka() {
	for {
		select {
		case msg := <-MsgChan:
			_, _, err := kafkaClient.SendMessage(msg)
			if err != nil {
				logger.Errorf("send msg to kafka failed, err:%v", err)
				return
			}
			logger.Infof("send msg to kafka success, msg:%v", msg)
		}
	}
}
