include ${PWD}/migrate/.env

init:
	docker-compose build && docker-compose up -d && docker-compose ps
up:
	docker-compose up -d
ps:
	docker-compose ps
logs.all:
	docker-compose logs -f
logs:
	docker-compose logs -f ${s}
exec:
	docker-compose exec ${s} sh
down:
	docker-compose down
migrate.create:
	docker-compose exec migrate migrate create -ext sql -dir migrations ${name}
migrate.up:
	docker-compose exec migrate migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):5432/$(DB_NAME)?sslmode=disable" -path migrations up
migrate.down:
	docker-compose exec migrate migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):5432/$(DB_NAME)?sslmode=disable" -path migrations down
migrate.up.test:
	docker-compose exec migrate migrate -database "postgres://$(DB_USER_TEST):$(DB_PASSWORD_TEST)@$(DB_HOST_TEST):5432/$(DB_NAME_TEST)?sslmode=disable" -path migrations up
migrate.down.test:
	docker-compose exec migrate migrate -database "postgres://$(DB_USER_TEST):$(DB_PASSWORD_TEST)@$(DB_HOST_TEST):5432/$(DB_NAME_TEST)?sslmode=disable" -path migrations down
