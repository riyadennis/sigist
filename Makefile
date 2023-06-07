
benthos:
	benthos create > config.yaml
docker-run:
	docker-compose -f environment/docker-compose.yaml up --build
docker-clean:
	docker-compose down --rmi all
