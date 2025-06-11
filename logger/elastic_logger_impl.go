package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/sirupsen/logrus"
)

var (
	elasticLoggerInstance ElasticLogger
	onceLogger            sync.Once
)

type ElasticLoggerImpl struct {
	logger *logrus.Logger
}

type ElasticHookImpl struct {
	client    *elasticsearch.Client
	indexName string
}

func NewElasticHookImpl(client *elasticsearch.Client, index string) ElasticHook {
	return &ElasticHookImpl{
		client:    client,
		indexName: index,
	}
}

func (h *ElasticHookImpl) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.FatalLevel,
	}
}

func (h *ElasticHookImpl) Fire(entry *logrus.Entry) error {
	doc := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"level":     entry.Level.String(),
		"message":   entry.Message,
		"fields":    entry.Data,
	}

	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	res, err := h.client.Index(
		h.indexName,
		bytes.NewReader(data),
		h.client.Index.WithContext(context.Background()),
		h.client.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elastic hook failed with status: %s", res.Status())
	}

	return nil
}

func NewElasticLoggerImpl(esClient *elasticsearch.Client, index string) ElasticLogger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.AddHook(NewElasticHookImpl(esClient, index))
	return &ElasticLoggerImpl{
		logger: log,
	}
}

func (l *ElasticLoggerImpl) Log(entity string, level string, message string) {
	entry := l.logger.WithFields(logrus.Fields{
		"entity": entity,
		"level":  level,
	})

	switch level {
	case "debug":
		entry.Debug(message)
	case "info":
		entry.Info(message)
	case "warn":
		entry.Warn(message)
	case "error":
		entry.Error(message)
	case "fatal":
		entry.Fatal(message)
	default:
		entry.Info(message)
	}
}

func GetLogger(index string) ElasticLogger {
	onceLogger.Do(func() {
		cfg := elasticsearch.Config{
			Addresses: []string{"http://localhost:9200"},
		}

		esClient, err := elasticsearch.NewClient(cfg)
		if err != nil {
			log.Fatal("FATAL: ", err)
		}

		elasticLoggerInstance = NewElasticLoggerImpl(esClient, index)
	})

	return elasticLoggerInstance
}
