package kafka

import (
	"context"
	"fmt"
	"testing"
)

func KafkaTest(t *testing.T) {
	client := Client{
		Servers:            []string{"0.0.0.0:9092"},
		WriteConnTimeoutMs: 5000,
		ReadConnTimeoutMs:  5000,
		GroupID:            "test",
	}
	//生成读取进程
	client.NewReadMessages("test")

	defer func() {
		err := client.ReadClose()
		if err != nil {
			fmt.Printf("failed to close reader: %s", err)
		}
	}()

	for {
		m, err := client.Reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Value))
	}
}
