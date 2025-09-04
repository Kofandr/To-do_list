.PHONE: run start help lint migrate

MAIN_PATH = cmd/app/main.go
DOCS_DIR = docs
MIGRATIONS_DIR = migrations

run:
	@env $$(grep -v '^#' .env | xargs) go run $(MAIN_PATH)

start: lint run
	go run cmd/app/main.go --migrate

lint:
	golangci-lint run

migrate:
	goose -dir $(MIGRATIONS_DIR) postgres $(DATABASE_URL) up

help:
	@echo "Доступные команды:"
	@echo "  make run     - Полный запуск (очистка + генерация + запуск)"
	@echo "  make start   - Быстрый запуск (без очистки и генерации)"
	@echo "  make lint    - Запустить линтеры"
	@echo "  make migrate - Применить миграции с помощью goose"
	@echo "  make help    - Показать эту справку"