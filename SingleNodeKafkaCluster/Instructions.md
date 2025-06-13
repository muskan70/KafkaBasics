### Commands To run Single node Kafka server

1. Start kafka server on docker
> docker compose up -d

2. Open a command terminal on the Kafka container
> docker exec -it -w /opt/kafka/bin {broker-name} sh

3. create a topic:
> ./kafka-topics.sh --create --topic {topic-name} --bootstrap-server {broker-listener->host:port}

4. Consume the messages with this command and use CTRL-C to close the consumer.
> ./kafka-console-consumer.sh --topic {topic-name} --from-beginning --bootstrap-server {broker-listener->host:port}

5. Start a console producer with this command and then push messages to kafka broker, use CTRL-C to close the producer.
> ./kafka-console-producer.sh  --topic {topic-name} --bootstrap-server {broker-listener->host:port}

6. To shut down the container, run
> docker compose down -v

Note: {broker-listener->host:port} = broker:29092