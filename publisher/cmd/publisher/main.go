package main

import (
	"context"
	"encoding/json"
	"github.com/AMCR/kafka-publisher/pkg/swagger-clients/archive/client/archive"
	"github.com/AMCR/kafka-publisher/pkg/swagger-clients/archive/models"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
)

const NYT_URL = "http://localhost:8080/svc/archive/v1/"

func sendMessageFnc(producer *kafka.Producer) func(ctx context.Context, correlationID uuid.UUID, topic string, key string, msg []byte) error {
	return func(ctx context.Context, correlationID uuid.UUID, topic string, key string, data []byte) error {
		msg := &kafka.Message{
			Headers: []kafka.Header{{
				Key: "x-correlation-id", Value: []byte(correlationID.String()),
			}},
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          data,
			Key:            []byte(key),
		}
		return producer.Produce(msg, nil)
	}
}

func run(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})
	if err != nil {
		return err
	}

	defer func() {
		producer.Flush(1000)
		producer.Close()
		wg.Wait()
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case e, ok := <-producer.Events():
				if !ok {
					return
				}
				switch msg := e.(type) {
				case *kafka.Message:
					topic := msg.TopicPartition
					if topic.Error != nil {
						slog.Error("fail to create message due to client error: %s", topic.Error)
					} else {
						slog.Debug("message delivered to topic %s>%s with offset [%s]", topic.Topic, topic.Partition, topic.Offset)
					}
				case *kafka.Error:
					slog.Warn("fail to create message due to client error: %s", msg.Error())
				}
			}
		}
	}()
	wg.Add(1)

	nytClient, err := createNYTClient()
	if err != nil {
		return err
	}

	sendMsgFnc := sendMessageFnc(producer)
	articles, err := consumeArticles(ctx, nytClient, 2024, 1)
	if err != nil {
		return err
	}

	for _, article := range articles {
		correlationID := uuid.New()
		payload, err := json.Marshal(article)
		if err != nil {
			return err
		}
		slog.Info("sending message", "correlation_id", correlationID.String())
		err = sendMsgFnc(ctx, correlationID, "articles", article.ID, payload)
		if err != nil {
			slog.Error("fail to send message with error: %s", err)
		} else {
			slog.Info("message send successfully")
		}
	}
	return nil
}

func createNYTClient() (archive.ClientService, error) {
	httpClient := &http.Client{}

	urlP, err := url.Parse(NYT_URL)
	if err != nil {
		return nil, err
	}
	cli := client.NewWithClient(
		urlP.Host,
		urlP.Path,
		[]string{urlP.Scheme},
		httpClient,
	)
	return archive.New(cli, strfmt.Default), nil
}

func consumeArticles(ctx context.Context, client archive.ClientService, year, month int) ([]*models.Article, error) {
	res, err := client.GetYearMonthJSON(archive.NewGetYearMonthJSONParams().
		WithContext(ctx).
		WithYear(int64(year)).
		WithMonth(int64(month)),
		nil,
	)
	if err != nil {
		return nil, err
	}
	if res == nil || res.Payload == nil || res.Payload.Response == nil {
		return nil, nil
	}
	return res.Payload.Response.Docs, nil
}

func main() {
	ctx, cnlFnc := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cnlFnc()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "exiting application with error %s", err)
	}
}
