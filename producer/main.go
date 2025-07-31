package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	// наружные порты, проброшенные из compose
	brokers := []string{"localhost:9094", "localhost:9095", "localhost:9096"}

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5
	cfg.Version = sarama.V3_4_0_0 // Kafka 4.0 → протокол 3.4

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		log.Fatalf("producer init: %v", err)
	}
	defer producer.Close()

	for i := 0; ; i++ {
		msg := fmt.Sprintf("kraft-%d", i)

		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: "demo-topic",
			// можно задать Key, если нужно упорядочивание по partition
			//Key:   sarama.StringEncoder(entity_id),
			Value: sarama.StringEncoder(msg),
		})
		if err != nil {
			log.Printf("send error: %v", err)
		} else {
			log.Printf("sent %q to partition %d @ offset %d", msg, partition, offset)
		}
		time.Sleep(5 * time.Second)
	}
}
