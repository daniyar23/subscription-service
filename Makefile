include .env
export

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONS_PATH=migrations

.PHONY: run up stop down migup migdown migcreate


# Запуск контейнеров в фоне
up:
	docker-compose up -d

# Остановка контейнеров
stop:
	docker-compose stop

# Полное удаление контейнеров
down:
	docker-compose down

# Применение всех миграций
migup:

# Откат последней миграции
migdown:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

# Посмотреть статус миграций
migstatus:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

# команда для создания новой миграции:
migcreate:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
