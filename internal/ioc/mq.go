package ioc

import (
	"go_project_template/internal/event"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

type MQConfig struct {
	Addrs []string
}

// InitSaramaClient 初始化Sarama客户端，用于连接到Kafka集群
func InitSaramaClient() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}

	var cfg Config

	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}

	scfg := sarama.NewConfig()
	// 配置生产者需要返回确认成功的消息
	scfg.Producer.Return.Successes = true

	client, err := sarama.NewClient(cfg.Addrs, scfg)
	if err != nil {
		panic(err)
	}

	return client
}

// InitSyncProducer 使用已有的Sarama客户端初始化同步生产者
func InitSyncProducer(c sarama.Client) sarama.SyncProducer {
	// 根据现有的客户端实例创建同步生产者
	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}

	return p
}

// InitConsumers 初始化并返回一个事件消费者
func InitConsumers() []event.Consumer {
	// 返回消费者切片
	return []event.Consumer{}
}
