// consumer/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"
)

type handler struct {
	name string
}

func (h handler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h handler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h handler) ConsumeClaim(sess sarama.ConsumerGroupSession, cl sarama.ConsumerGroupClaim) error {
	fmt.Println("ConsumeClaim running...")
	for msg := range cl.Messages() {
		log.Printf("%s ⇠ partition=%d offset=%d value=%s",
			h.name, msg.Partition, msg.Offset, string(msg.Value))
		sess.MarkMessage(msg, "") // фиксируем оффсет для _этой_ группы
	}
	return nil
}

//обработка пачкой (не тестил)
//func (h handler) ConsumeClaim(sess sarama.ConsumerGroupSession,
//	cl sarama.ConsumerGroupClaim) error {
//
//	const (
//		batchSize   = 100                    // сколько сообщений в пачке
//		flushPeriod = 3 * time.Second        // или по таймеру
//	)
//	topic     := cl.Topic()
//	partition := cl.Partition()
//
//	buf   := make([]*sarama.ConsumerMessage, 0, batchSize)
//	timer := time.NewTimer(flushPeriod)
//	defer timer.Stop()
//
//	flush := func() {
//		if len(buf) == 0 {
//			return
//		}
//		// 1. бизнес-обработка
//		if err := h.insertBulk(buf); err != nil {
//			log.Printf("bulk err: %v", err)
//			return            // сообщения НЕ помечаем, вернутся заново
//		}
//		// 2. помечаем последний offset (+1)
//		last := buf[len(buf)-1]
//		sess.MarkOffset(topic, partition, last.Offset+1, "")
//		buf = buf[:0]         // очистили буфер
//		timer.Reset(flushPeriod)
//	}
//
//	for {
//		select {
//		case msg, ok := <-cl.Messages():
//			if !ok {
//				flush()      // закрытие claim'а перед ребалансом
//				return nil   // выход из ConsumeClaim
//			}
//			buf = append(buf, msg)
//			if len(buf) >= batchSize {
//				flush()
//			}
//
//		case <-timer.C:
//			flush()
//		}
//	}
//}

func main() {
	sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)

	groupID := os.Getenv("GROUP")
	if groupID == "" {
		groupID = "default"
		//log.Fatal("set GROUP env var, e.g.  GROUP=cg-test")
	}

	cfg := sarama.NewConfig()
	//cfg.Version = sarama.MaxVersion
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Version = sarama.V3_4_0_0
	cfg.Consumer.Group.Rebalance.Timeout = 60 * time.Second

	brokers := []string{"localhost:9094", "localhost:9095", "localhost:9096"}
	cg, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		log.Fatalf("consumer init: %v", err)
	}
	defer cg.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			if err := cg.Consume(ctx, []string{"demo-topic"}, handler{name: groupID}); err != nil {
				log.Printf("consume err: %v", err)
			}
			time.Sleep(5 * time.Second)
			// небольшая пауза, чтобы не крутить цикл на ошибке
		}
	}()

	fmt.Println("running...")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	cancel()
}
