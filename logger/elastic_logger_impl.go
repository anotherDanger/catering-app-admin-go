package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/sirupsen/logrus"
)

var (
	elasticLoggerInstance ElasticLogger
	isLoggerInitialized   bool
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
	if !isLoggerInitialized {
		cfg := elasticsearch.Config{
			Addresses: []string{"http://localhost:9200"},
		}

		esClient, err := elasticsearch.NewClient(cfg)
		if err != nil {
			logToFile("log/elasticsearchfatal.log", fmt.Sprintf("Failed to create Elasticsearch client: %v", err))
			elasticLoggerInstance = NewFileFallbackLogger("log/elasticsearchfatal.log")
			isLoggerInitialized = true
			return elasticLoggerInstance
		}

		res, err := esClient.Info()
		if err != nil {
			logToFile("log/elasticsearchfatal.log", fmt.Sprintf("Failed to connect to Elasticsearch: %v", err))
			elasticLoggerInstance = NewFileFallbackLogger("log/elasticsearchfatal.log")
			isLoggerInitialized = true
			if res != nil {
				res.Body.Close()
			}
			return elasticLoggerInstance
		}
		defer res.Body.Close()

		if res.IsError() {
			logToFile("log/elasticsearchfatal.log", fmt.Sprintf("Elasticsearch returned an error: %s", res.Status()))
			elasticLoggerInstance = NewFileFallbackLogger("log/elasticsearchfatal.log")
			isLoggerInitialized = true
			return elasticLoggerInstance
		}

		elasticLoggerInstance = NewElasticLoggerImpl(esClient, index)
		isLoggerInitialized = true
	}

	return elasticLoggerInstance
}

func logToFile(filePath, message string) {
	dir := filepath.Dir(filePath)
	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Printf("Failed to create log directory %s: %v\n", dir, err)
				filePath = filepath.Base(filePath)
			}
		}
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s - %s\n", time.Now().Format(time.RFC3339), message))
	if err != nil {
		fmt.Printf("Failed to write to log file: %v\n", err)
	}
}

func NewFileFallbackLogger(filePath string) ElasticLogger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	dir := filepath.Dir(filePath)
	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create fallback log file %s: %v. Logging to stdout instead.\n", dir, err)
				log.SetOutput(os.Stdout)
				return &ElasticLoggerImpl{logger: log}
			}
		}
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open fallback log file %s: %v. Logging to stdout instead.\n", filePath, err)
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(file)
	}

	return &ElasticLoggerImpl{logger: log}
}
