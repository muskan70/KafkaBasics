docker exec -it -w /opt/kafka/bin kafka1 sh 
./kafka-topics.sh --create --topic inventory-ticks --partitions 3 --replication-factor 2 --bootstrap-server kafka1:29092