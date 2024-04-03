.PHONY: all up migrate generator data_processing fmt

all: up migrate

up:
	docker-compose up -d

migrate:
	docker-compose exec database sh -c 'psql -U casino < /db/migrations/00001.create_base.sql'

data_processing:
	docker-compose run --rm -p 8080:8080 data_processing

generator:
	docker-compose run --rm generator

fmt:
	go fmt ./...