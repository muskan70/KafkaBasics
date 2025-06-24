docker exec -it -w /opt/kafka/bin kafka1 sh 
./kafka-console-producer.sh  --topic inventory-ticks --bootstrap-server kafka1:29092