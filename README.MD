# Kafka Cluster with Golang Producer & Consumer (KRaft Mode)

This project sets up a multi-node Apache Kafka cluster in **KRaft (Kafka Raft Metadata mode)** using Docker Compose. It includes **three Kafka controllers** and **three brokers**, configured with replication and external access. The project also provides **Golang-based producer and consumer clients** for testing message production and consumption.

## 📦 Project Structure
├── kafka/docker-compose.yml # Kafka + Controller cluster configuration
├── producer/ # Golang Kafka producer example
└── consumers/ # Golang Kafka consumer group example


## 🧩 Features

- Kafka 4.0.0 running in **KRaft mode** (no Zookeeper)
- 3 controller nodes for metadata quorum
- 3 broker nodes with external access (`9094`, `9095`, `9096`)
- Golang sample producer and consumer (Sarama-based) for testing
- Configurable via environment anchors for simplicity
- Auto-persistent volumes for Kafka logs

## 🛠 Requirements

- Docker
- Docker Compose
- Go 1.18+ (for running Golang clients)

## 🚀 How to Start the Cluster

# Шаг 0 - удаление топика (опционально). Step 0 - delete the topic (optional)
```bash
docker compose exec broker1 \
bash -c "/opt/kafka/bin/kafka-topics.sh \
--delete \
--topic demo-topic \
--bootstrap-server broker1:9092"
```

# Шаг 1 — генерируем ID (один раз) oIN7jCltQFaDgQWgDZKuxA. Step 1 - gen uuid for cluster
```bash
export KAFKA_CLUSTER_ID=$(uuidgen | tr 'A-Z' 'a-z')
# Подставьте значение в docker-compose.yml вместо oIN7jClt...
```

#  Шаг 2 — Поднимаем кластер. Step 2 docker up
```bash
docker compose up -d               # ~15-20 s на прогрев JDK
docker compose logs -f broker1     # ждём строку "Kafka startTimeMs"
```
# Шаг 3 — Создаём продовый топик. Step 3 create a topic close to prod version
```bash
docker compose exec broker1 \
  bash -c "/opt/kafka/bin/kafka-topics.sh \
           --create \
           --topic demo-topic \
           --bootstrap-server broker1:9092 \
           --partitions 1 \
           --replication-factor 3"
```

# Шаг 4 — Проверяем распределение реплик. Step 4 - check the replicas
```bash
docker compose exec broker1 \
  bash -c "/opt/kafka/bin/kafka-topics.sh \
             --describe \
             --topic demo-topic \
             --bootstrap-server broker1:9092"

```
# Шаг 5 — Тестовая отправка / чтение. Step 5 test read|write

## Producer
```bash
docker compose exec broker2 \
  bash -c "/opt/kafka/bin/kafka-console-producer.sh \
           --bootstrap-server broker2:9092 \
           --topic demo-topic"
```

## Consumer
```bash
docker compose exec broker3 \
  bash -c "/opt/kafka/bin/kafka-console-consumer.sh \
           --bootstrap-server broker3:9092 \
           --topic demo-topic \
           --from-beginning"
```
