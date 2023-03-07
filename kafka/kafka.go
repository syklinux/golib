package kafka

import (
	"context"
	"time"

	"github.com/syklinux/golib/log"

	"github.com/segmentio/kafka-go"
)

type Client struct {
	Servers            []string `json:"Servers"`
	WriteConnTimeoutMs int      `json:"WriteConnTimeoutMs"`
	ReadConnTimeoutMs  int      `json:"ReadConnTimeoutMs"`
	GroupID            string   `json:"GroupID"`
	DefaultTopic       string   `json:"DefaultTopic"`
	Reader             *kafka.Reader
}

type Message struct {
	Topic string `json:"topic"`
	Msg   string `json:"msg"`
	Key   string `json:"key"`
}

// ProduceMessage 发送一条消息
func (client *Client) ProduceMessage(ctx context.Context, message Message) error {
	var (
		err error
	)

	write := &kafka.Writer{
		Addr:         kafka.TCP(client.Servers...),
		Balancer:     &kafka.LeastBytes{},
		Topic:        message.Topic,
		WriteTimeout: time.Millisecond * time.Duration(client.WriteConnTimeoutMs),
		ReadTimeout:  time.Millisecond * time.Duration(client.ReadConnTimeoutMs),
	}

	err = write.WriteMessages(ctx, kafka.Message{
		Value: []byte(message.Msg),
		Key:   []byte(message.Key),
	})

	if err != nil {
		log.Errorf("ProduceMessage failed to write message:", err)
		return err
	}

	if err = write.Close(); err != nil {
		log.Errorf("ProduceMessage failed to close:", err)
		return err
	}

	return nil
}

// ProduceMessages 发送多条消息, 当topic不为空时，message不能包含topic字段
func (client *Client) ProduceMessages(ctx context.Context, topic string, messages []Message) (err error) {

	write := &kafka.Writer{
		Addr:         kafka.TCP(client.Servers...),
		Balancer:     &kafka.LeastBytes{},
		Topic:        topic,
		WriteTimeout: time.Millisecond * time.Duration(client.WriteConnTimeoutMs),
		ReadTimeout:  time.Millisecond * time.Duration(client.ReadConnTimeoutMs),
	}

	msgs := make([]kafka.Message, len(messages))

	for _, v := range messages {
		var msg kafka.Message
		msg.Value = []byte(v.Msg)
		msg.Key = []byte(v.Key)
		msg.Topic = v.Topic
		msgs = append(msgs, msg)
	}
	err = write.WriteMessages(ctx, msgs...)

	if err != nil {
		log.Errorf("ProduceDiffTopicMessages failed to write message:", err)
		return err
	}

	if err = write.Close(); err != nil {
		log.Errorf("ProduceDiffTopicMessages failed to close:", err)
		return err
	}

	return nil
}

// NewReadMessages 生成读取进程，需要自己close掉
func (client *Client) NewReadMessages(topic string) {
	read := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  client.Servers,
		GroupID:  client.GroupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	client.Reader = read
}

// ReadClose 读取进程关闭
func (client *Client) ReadClose() (err error) {
	err = client.Reader.Close()
	if err != nil {
		return err
	}
	return nil
}
