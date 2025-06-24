docker exec -it -w /opt/kafka/bin kafka1 sh 
./kafka-console-consumer.sh --topic inventory-ticks --from-beginning --bootstrap-server kafka1:29092 --group group1