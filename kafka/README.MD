# Apache Kafka 4.0.0 (KRaft) cluster 3×3

## Требования
* Docker 24+
* Docker Compose v2.20+

# Шаг 1 — Клонируем репозиторий и генерируем ID (один раз) oIN7jCltQFaDgQWgDZKuxA
```bash
git clone <repo> kafka-kraft-prod
cd kafka-kraft-prod
export KAFKA_CLUSTER_ID=$(uuidgen | tr 'A-Z' 'a-z')
# Подставьте значение в docker-compose.yml вместо oIN7jClt...
```

#  Шаг 2 — Поднимаем кластер
```
docker compose up -d               # ~15-20 s на прогрев JDK
docker compose logs -f broker1     # ждём строку "Kafka startTimeMs"
```
# Шаг 3 — Создаём продовый топик
```
docker compose exec broker1 \
  bash -c "/opt/kafka/bin/kafka-topics.sh \
           --create \
           --topic demo-topic \
           --bootstrap-server broker1:9092 \
           --partitions 6 \
           --replication-factor 3"
```

# Шаг 4 — Проверяем распределение реплик
```
docker compose exec broker1 \
  bash -c "/opt/kafka/bin/kafka-topics.sh \
             --describe \
             --topic demo-topic \
             --bootstrap-server broker1:9092"

```
# Шаг 5 — Тестовая отправка / чтение

## Producer
```
docker compose exec broker2 \
  bash -c "/opt/kafka/bin/kafka-console-producer.sh \
           --bootstrap-server broker2:9092 \
           --topic demo-topic"
```

## Consumer
```
docker compose exec broker3 \
  bash -c "/opt/kafka/bin/kafka-console-consumer.sh \
           --bootstrap-server broker3:9092 \
           --topic demo-topic \
           --from-beginning"
```

добавить новые brokerN с уникальным KAFKA_CFG_NODE_ID и тем же KAFKA_CLUSTER_ID).
::contentReference[oaicite:1]{index=1}
