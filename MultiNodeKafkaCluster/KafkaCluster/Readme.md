### Commands To run Multi node Kafka server

1. Create Volumes
> docker volume create kafka1_data<br>
> docker volume create kafka2_data<br>
> docker volume create kafka3_data<br>
> docker run --rm -v kafka1_data:/data alpine chown -R 1001:1001 /data<br>
> docker run --rm -v kafka2_data:/data alpine chown -R 1001:1001 /data<br>
> docker run --rm -v kafka3_data:/data alpine chown -R 1001:1001 /data<br>

2. Start kafka server on docker : 3 nodes
> docker compose up -d

3. Open a command terminal on the Kafka container
> docker exec -it -w /opt/kafka/bin {broker-name} sh

4. create a topic with 3 partitions
> ./kafka-topics.sh --create --topic {topic-name} --partitions 3 --replication-factor 2 --bootstrap-server {broker-listener->host:port}

5. Start a console producer with this command and then push messages to kafka broker, use CTRL-C to close the producer.
> ./kafka-console-producer.sh  --topic {topic-name} --bootstrap-server {broker-listener->host:port}

6. Consume the messages with this command and use CTRL-C to close the consumer : open on 2 tabs to get 2 consumers with groop name
> ./kafka-console-consumer.sh --topic {topic-name} --from-beginning --bootstrap-server {broker-listener->host:port} --group {group-name}

7. To shut down the container, run
> docker compose down -v

Note: 
1. {broker-listener->host:port} = kafka1:29092
2. To check topic partititons:
> ./kafka-topics.sh --describe --topic mytopic --bootstrap-server kafka1:29092
3. To check consumer group status
> ./kafka-consumer-groups.sh --bootstrap-server kafka1:29092 --group group1 --describe
