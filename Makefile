
benthos:
	benthos create > config.yaml
docker-run:
	docker-compose -f environment/docker-compose.yaml up --build
docker-clean:
	docker-compose -f environment/docker-compose.yaml down --rmi all
check-logs:
    docker-compose -f environment/docker-compose.yaml exec kafka cat /var/log/broker.log 2>&1 | tail -100
ssh-rest-service:
    docker-compose -f environment/docker-compose.yaml exec rest-service /bin/bash
