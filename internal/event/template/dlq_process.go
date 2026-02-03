package template

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const (
	GroupTemplateDLQ  = "template_dlq"
	DLQProcessTimeout = 10 * time.Second
	DLQMaxRetries     = 5
	DLQBaseWaitTime   = 5 * time.Second
)

type TemplateDeadLetterConsumer struct {
	client sarama.Client
	logger *zap.Logger
}

func NewTemplateDeadLetterConsumer(
	client sarama.Client,
	logger *zap.Logger,
) *TemplateDeadLetterConsumer {
	return &TemplateDeadLetterConsumer{
		client: client,
		logger: logger,
	}
}

func (p *TemplateDeadLetterConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient(GroupTemplateDLQ, p.client)
	if err != nil {
		p.logger.Error("创建死信队列消费者组失败", zap.Error(err))
		return err
	}

	p.logger.Info("DeadLetterConsumer 开始消费死信队列")

	// 监听消费者组错误，避免吞掉底层错误信息
	go func() {
		for err := range cg.Errors() {
			p.logger.Error("死信队列消费者组错误", zap.Error(err))
		}
	}()

	// 启动死信队列消息消费
	go func() {
		defer cg.Close()
		for {
			select {
			case <-ctx.Done():
				p.logger.Info("死信队列消费者停止")
				return
			default:
				if err := cg.Consume(ctx, []string{TopicDeadLetter}, &dlqConsumerGroupHandler{consumer: p}); err != nil {
					p.logger.Error("死信队列消费循环出错", zap.Error(err))
					time.Sleep(500 * time.Millisecond)
				}
			}
		}
	}()

	return nil
}

type dlqConsumerGroupHandler struct {
	consumer *TemplateDeadLetterConsumer
}

func (h *dlqConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *dlqConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *dlqConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var err error
		for i := 0; i < DLQMaxRetries; i++ {
			if err = h.consumer.processDLQMessage(msg); err == nil {
				break
			}

			if i < DLQMaxRetries-1 { // 最后一次失败不需要记录重试日志
				h.consumer.logger.Error("处理死信消息失败,准备重试",
					zap.Error(err),
					zap.Int("重试次数", i+1),
					zap.Int("剩余重试次数", DLQMaxRetries-i-1),
					zap.ByteString("message", msg.Value))

				// 指数退避策略,等待时间随重试次数指数增长
				waitTime := DLQBaseWaitTime * time.Duration(1<<uint(i))
				select {
				case <-time.After(waitTime):
				case <-sess.Context().Done():
					return sess.Context().Err()
				}
			}
		}

		if err != nil {
			h.consumer.logger.Error("处理死信消息最终失败",
				zap.Error(err),
				zap.ByteString("message", msg.Value))
		} else {
			h.consumer.logger.Info("死信消息处理成功",
				zap.ByteString("message", msg.Value))
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

// processDLQMessage 处理死信队列中的消息
func (p *TemplateDeadLetterConsumer) processDLQMessage(msg *sarama.ConsumerMessage) error {
	var evt TemplateEvent

	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		p.logger.Error("死信消息反序列化失败", zap.Error(err), zap.ByteString("message", msg.Value))
		return fmt.Errorf("死信消息反序列化失败: %w", err)
	}

	// 从死信队列获取原始主题、时间等信息
	originalTopic := ""
	for _, header := range msg.Headers {
		if string(header.Key) == "original_topic" {
			originalTopic = string(header.Value)
		}
	}

	p.logger.Info("处理死信消息",
		zap.String("original_topic", originalTopic),
		zap.ByteString("message", msg.Value))

	ctx, cancel := context.WithTimeout(context.Background(), DLQProcessTimeout)
	defer cancel()

	// 重新处理消息
	if err := p.handleDeadLetterMessage(ctx, &evt); err != nil {
		p.logger.Error("处理死信消息失败", zap.Error(err))
		return err
	}

	return nil
}

// handleDeadLetterMessage 处理死信消息的具体业务逻辑
func (p *TemplateDeadLetterConsumer) handleDeadLetterMessage(ctx context.Context, evt *TemplateEvent) error {
	// ...
	p.logger.Info("成功处理死信消息", zap.Int64("template_id", evt.TemplateId))
	return nil
}
