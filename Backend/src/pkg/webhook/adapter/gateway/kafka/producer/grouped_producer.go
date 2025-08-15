package producer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/segmentio/kafka-go"
)

type Message struct {
	Key   string
	Value interface{}
}

type GroupedProducer struct {
	writer   *kafka.Writer
	channels map[string]chan Message
	mu       sync.RWMutex
	workers  int
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	log      logging.Logger
}

func NewGroupedProducer(brokers []string, topic string, workers int, log logging.Logger) *GroupedProducer {
	ctx, cancel := context.WithCancel(context.Background())
	log.Info("Initializing Kafka grouped producer", map[string]interface{}{
		"brokers": brokers,
		"topic":   topic,
		"workers": workers,
	})

	return &GroupedProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.Hash{},
			Async:    true,
		},
		channels: make(map[string]chan Message),
		workers:  workers,
		ctx:      ctx,
		cancel:   cancel,
		log:      log,
	}
}

func (p *GroupedProducer) Start() {
	p.log.Info("Starting Kafka producer workers", map[string]interface{}{
		"worker_count": p.workers,
	})

	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		workerID := i + 1
		p.log.Debug("Starting worker", map[string]interface{}{
			"worker_id": workerID,
		})
		go p.worker(workerID)
	}
}

func (p *GroupedProducer) Stop() {
	p.log.Info("Stopping Kafka producer", nil)
	p.cancel()
	p.wg.Wait()

	p.mu.Lock()
	defer p.mu.Unlock()

	p.log.Debug("Closing producer channels", map[string]interface{}{
		"channel_count": len(p.channels),
	})

	for key, ch := range p.channels {
		p.log.Debug("Closing channel", map[string]interface{}{
			"key": key,
		})
		close(ch)
	}
	p.log.Info("Kafka producer stopped successfully", nil)
}

func (p *GroupedProducer) Produce(key string, value interface{}) {
	p.log.Debug("Attempting to produce message", map[string]interface{}{
		"key": key,
	})

	p.mu.Lock()
	ch, exists := p.channels[key]
	if !exists {
		p.log.Debug("Creating new channel for key", map[string]interface{}{
			"key": key,
		})
		ch = make(chan Message, 100)
		p.channels[key] = ch
	}
	p.mu.Unlock()

	// Exponential backoff parameters
	maxRetries := 5
	initialBackoff := 100 * time.Millisecond
	maxBackoff := 5 * time.Second
	currentBackoff := initialBackoff

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case ch <- Message{Key: key, Value: value}:
			p.log.Debug("Message queued successfully", map[string]interface{}{
				"key": key,
			})
			return
		case <-time.After(currentBackoff):
			p.log.Warn("Channel full, retrying with backoff",
				map[string]interface{}{
					"key":            key,
					"retry":          retry + 1,
					"backoff_ms":     currentBackoff.Milliseconds(),
					"max_retries":    maxRetries,
					"channel_length": len(ch),
				})

			currentBackoff = time.Duration(float64(currentBackoff) * 1.5)
			if currentBackoff > maxBackoff {
				currentBackoff = maxBackoff
			}
		}
	}

	p.log.Error("Failed to send message after max retries",
		map[string]interface{}{
			"key":            key,
			"max_retries":    maxRetries,
			"channel_length": len(ch),
		})
}

func (p *GroupedProducer) worker(workerID int) {
	defer p.wg.Done()

	p.log.Info("Worker started", map[string]interface{}{
		"worker_id": workerID,
	})

	messageCount := 0
	lastLog := time.Now()
	logInterval := 5 * time.Second

	for {
		select {
		case <-p.ctx.Done():
			p.log.Info("Worker shutting down", map[string]interface{}{
				"worker_id":          workerID,
				"messages_processed": messageCount,
			})
			return
		default:
			p.mu.RLock()
			channels := make([]chan Message, 0, len(p.channels))
			for _, ch := range p.channels {
				channels = append(channels, ch)
			}
			p.mu.RUnlock()

			for _, ch := range channels {
				select {
				case msg := <-ch:
					valBytes, ok := msg.Value.([]byte)
					if !ok {
						p.log.Error("Invalid message value type", map[string]interface{}{
							"worker_id": workerID,
							"key":       msg.Key,
							"type":      fmt.Sprintf("%T", msg.Value),
						})
						continue
					}

					err := p.writer.WriteMessages(p.ctx, kafka.Message{
						Key:   []byte(msg.Key),
						Value: valBytes,
					})
					if err != nil {
						p.log.Error("Failed to write message to Kafka", map[string]interface{}{
							"worker_id": workerID,
							"key":       msg.Key,
							"error":     err.Error(),
						})
						continue
					}

					messageCount++
					if time.Since(lastLog) > logInterval {
						p.log.Info("Worker status", map[string]interface{}{
							"worker_id":          workerID,
							"messages_processed": messageCount,
							"messages_per_sec":   float64(messageCount) / time.Since(lastLog).Seconds(),
						})
						messageCount = 0
						lastLog = time.Now()
					}

				default:
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
