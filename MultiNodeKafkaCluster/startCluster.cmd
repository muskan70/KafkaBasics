docker volume create kafka1_data
docker volume create kafka2_data
docker volume create kafka3_data
docker run --rm -v kafka1_data:/data alpine chown -R 1001:1001 /data
docker run --rm -v kafka2_data:/data alpine chown -R 1001:1001 /data
docker run --rm -v kafka3_data:/data alpine chown -R 1001:1001 /data
docker compose up -d