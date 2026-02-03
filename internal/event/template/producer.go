package template

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const TopicTemplateEvent = "template_events"

type Producer interface {
	ProduceTemplateEvent(evt TemplateEvent) error
}

type TemplateEvent struct {
	TemplateId int64 `json:"template_id"`
}

type TemplateSaramaSyncProducer struct {
	producer sarama.SyncProducer
	logger   *zap.Logger
}

func NewTemplateSaramaSyncProducer(producer sarama.SyncProducer, logger *zap.Logger) Producer {
	return &TemplateSaramaSyncProducer{
		producer: producer,
		logger:   logger,
	}
}

func (s *TemplateSaramaSyncProducer) ProduceTemplateEvent(evt TemplateEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		s.logger.Error("Failed to marshal template event", zap.Error(err))
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: TopicTemplateEvent,
		Value: sarama.StringEncoder(val),
	}

	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		s.logger.Error("Failed to send template event message", zap.Error(err))
		return err
	}

	s.logger.Info("Template event message sent", zap.Int32("partition", partition), zap.Int64("offset", offset))
	return nil
}
