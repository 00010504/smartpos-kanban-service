package kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/Invan2/invan_kanban_service/pkg/logger"
	"github.com/Shopify/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// Publisher ...
type Publisher struct {
	topic            string
	cloudEventClient cloudevents.Client
	sender           *kafka_sarama.Sender
}

func (kafka *Kafka) AddPublisher(topic string) {
	if kafka.publishers[topic] != nil {
		kafka.log.Warn("publisher exists", logger.Error(errors.New("publisher with the same topic already exists: "+topic)))
		return
	}

	sender, err := kafka_sarama.NewSender(
		[]string{kafka.cfg.KafkaUrl}, // Kafka connection url
		kafka.saramaConfig,           // Kafka sarama config
		topic,                        // Topic
	)

	if err != nil {
		panic(err)
	}

	// defer sender.Close(context.Background())

	c, err := cloudevents.NewClient(sender, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		panic(err)
	}

	kafka.publishers[topic] = &Publisher{
		topic:            topic,
		cloudEventClient: c,
		sender:           sender,
	}
}

// Push ...
func (r *Kafka) Push(topic string, e cloudevents.Event) error {
	p := r.publishers[topic]

	if p == nil {
		return fmt.Errorf("publisher with that topic doesn't exists: %s", topic)
	}

	result := p.cloudEventClient.Send(
		kafka_sarama.WithMessageKey(context.Background(), sarama.StringEncoder(e.ID())),
		e,
	)

	if cloudevents.IsUndelivered(result) {
		return fmt.Errorf("failed to publish event, Error: %v", result)
	}

	return nil
}
