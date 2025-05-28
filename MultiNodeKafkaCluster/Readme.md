### Commands To run Multi node Kafka server

1. Start kafka server on docker : 3 nodes
> docker compose up -d

2. Open a command terminal on the Kafka container
> docker exec -it -w /opt/kafka/bin {broker-name} sh

3. create a topic with 3 partitions
> ./kafka-topics.sh --create --topic {topic-name} --partitions 3 --bootstrap-server broker:29092

4. Start a console producer with this command and then push messages to kafka broker, use CTRL-C to close the producer.
> ./kafka-console-producer.sh  --topic {topic-name} --bootstrap-server broker:29092

5. Consume the messages with this command and use CTRL-C to close the consumer : open on 2 tabs to get 2 consumers with groop name
> ./kafka-console-consumer.sh --topic {topic-name} --from-beginning --bootstrap-server broker:29092 --group {group-name}

6. To shut down the container, run
> docker compose down -v