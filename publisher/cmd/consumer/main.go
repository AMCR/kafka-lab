package main

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

func run(ctx context.Context) error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "smallest",
	})
	if err != nil {
		return err
	}

	defer func() {
		_ = consumer.Unsubscribe()
		_ = consumer.Close()
	}()

	if err := consumer.Subscribe("articles", nil); err != nil {
		return err
	}

	kTopic := "articles"
	res, err := consumer.Position([]kafka.TopicPartition{
		{
			Topic: &kTopic,
		},
	})
	if err != nil {
		return err
	}
	_ = res
	slog.Info("start consuming messages:")
	for {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		default:
			event := consumer.Poll(100)
			switch msg := event.(type) {
			case *kafka.Message:
				attrs := []any{
					slog.String("Key", string(msg.Key)),
					slog.String("Timestamp", msg.Timestamp.Format(time.RFC3339Nano)),
					slog.String("Topic.Partition", *msg.TopicPartition.Topic),
					slog.String("Topic.Offset", msg.TopicPartition.Offset.String()),
				}
				for _, header := range msg.Headers {
					attrs = append(attrs, slog.String(fmt.Sprintf("Header.%s", header.Key), string(header.Value)))
				}
				slog.Info("message received", attrs...)
				_, err = consumer.CommitMessage(msg)
				if err != nil {
					slog.Error("fail to commit message", "error", err)
				}
			case kafka.Error:
				slog.Error("error received", "error", msg.Error())
			case kafka.PartitionEOF:
				slog.Info("reached eof partition", "msg", event.String())
			default:
				if event != nil {
					slog.Info("ignored message", "message", event.String())
				}
			}
		}
	}
}

func main() {
	ctx, cnl := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cnl()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "fail to run command", "error", err)
	}
}
