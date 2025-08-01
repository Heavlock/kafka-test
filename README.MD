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

```bash
docker compose up -d
