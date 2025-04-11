run:
	docker-compose up -d

build:
	docker-compose build

migrate-add: ## Create new migration file, usage: migrate-add [name=<migration_name>]
	goose -dir database/migrations create $(name) sql