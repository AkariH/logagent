package service

import (
	"github.com/IBM/sarama"
	"logagent/config"
	"logagent/kafka"
	"logagent/tailfile"
	"logagent/utils/logger"
	"time"
)

func Run() error {

	for {
		line, ok := <-tailfile.Tailfile.Lines
		if !ok {
			logger.Warnf("tail file close reopen, filename:%s\n", tailfile.Tailfile.Filename)
			time.Sleep(time.Second)
			continue
		}
		logger.Infof("msg is %s", line.Text)
		msg := &sarama.ProducerMessage{}
		cfg := config.GetConfig()
		msg.Topic = cfg.KafkaConfig.Topic
		msg.Value = sarama.StringEncoder(line.Text)

		kafka.MsgChan <- msg
	}

	return nil
}
