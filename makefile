-include src/.env

.EXPORT_ALL_VARIABLES:
	APP_HOST=${APP_HOST}
	APP_PORT=${APP_PORT}
	POSTGRES_USER=${POSTGRES_USER}
	POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
	POSTGRES_DB=${POSTGRES_DB}

init_docs:
	cd src && /home/vboxuser/go/bin/swag init

run_local:
	cd src && go run main.go

create_directories:
	mkdir -p graphite
	mkdir -p data

init:
	docker swarm init || true

setup: init create_directories

build:
	docker rmi http_go -f || true
	docker build -t http_go .

remove:
	if docker stack ls | grep -q http_go; then \
		docker stack rm http_go; \
	fi

deploy: remove build
	until \
	docker stack deploy \
	-c docker-compose.yaml \
	http_go; \
	do sleep 1; \
	done

logs:
	docker service logs http_go_app -f