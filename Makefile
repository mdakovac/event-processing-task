.PHONY: all up migrate generator data_processing fmt

all: up migrate

up:
	docker-compose up -d

migrate:
	docker-compose exec database sh -c 'psql -U casino < /db/migrations/00001.create_base.sql'

data_processing:
	docker-compose run --rm data_processing

generator:
	docker-compose run --rm generator

fmt:
	go fmt ./...