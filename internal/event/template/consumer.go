package template

import (
	"context"
	"encoding/json"
	"fmt"
	"go_project_template/internal/repository"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const (
	TopicDeadLetter = "template_events_dlq" // 死信队列主题
	MaxRetries      = 3                     // 最大重试次数
)

type TemplateEventConsumer struct {
	repo    repository.TemplateRepository
	client  sarama.Client
	logger  *zap.Logger
	dlqProd sarama.SyncProducer // 死信队列生产者
}

type consumerGroupHandler struct {
	consumer *TemplateEventConsumer
}

func NewTemplateEventConsumer(repo repository.TemplateRepository, client sarama.Client, dlqProd sarama.SyncProducer, logger *zap.Logger) *TemplateEventConsumer {
	return &TemplateEventConsumer{
		repo:    repo,
		client:  client,
		logger:  logger,
		dlqProd: dlqProd,
	}
}

// Start 启动消费者，并开始消费 Kafka 中的消息
func (p *TemplateEventConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("template_event", p.client)
	if err != nil {
		p.logger.Error("创建消费者组失败", zap.Error(err))
		return err
	}

	p.logger.Info("TemplateConsumer 开始消费")

	go func() {
		defer cg.Close()
		for {
			select {
			case <-ctx.Done():
				p.logger.Info("消费者停止")
				return
			default:
				// 开始消费指定的 Kafka 主题
				if err := cg.Consume(ctx, []string{TopicTemplateEvent}, &consumerGroupHandler{consumer: p}); err != nil {
					p.logger.Error("消费循环出错", zap.Error(err))
					continue
				}
			}
		}
	}()

	return nil
}

func (c *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// 处理每一条消息
		if err := c.consumer.processMessage(msg); err != nil {
			c.consumer.logger.Error("处理消息失败", zap.Error(err), zap.ByteString("message", msg.Value))
			// 发送到死信队列
			if err := c.consumer.sendToDLQ(msg); err != nil {
				c.consumer.logger.Error("发送到死信队列失败", zap.Error(err))
			}
			continue // 确保继续处理下一条消息
		}
		sess.MarkMessage(msg, "") // 只有成功处理后才标记
	}

	return nil
}

// sendToDLQ 发送消息到死信队列
func (p *TemplateEventConsumer) sendToDLQ(msg *sarama.ConsumerMessage) error {
	dlqMsg := &sarama.ProducerMessage{
		Topic: TopicDeadLetter,
		Value: sarama.ByteEncoder(msg.Value),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("original_topic"),
				Value: []byte(msg.Topic),
			},
			{
				Key:   []byte("error_time"),
				Value: []byte(time.Now().Format(time.RFC3339)),
			},
		},
	}

	_, _, err := p.dlqProd.SendMessage(dlqMsg)
	return err
}

// processMessage 处理从 Kafka 消费的消息
func (p *TemplateEventConsumer) processMessage(msg *sarama.ConsumerMessage) error {
	var event TemplateEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		p.logger.Error("反序列化消息失败",
			zap.Error(err),
			zap.String("message", string(msg.Value)))
		return fmt.Errorf("反序列化消息失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 处理消息
	if err := p.handleEvent(ctx, &event); err != nil {
		return err
	}

	return nil
}

// handleEvent 处理发布事件
func (p *TemplateEventConsumer) handleEvent(ctx context.Context, event *TemplateEvent) error {
	// ...
	return nil
}
