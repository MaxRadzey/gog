# Генерация Swagger-доки
swag:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/shortener/main.go -d . --parseInternal -o docs

# Форматирование кода + сортировка импортов
fmt:
	go run golang.org/x/tools/cmd/goimports@latest -w .

# Запуск приложения локально через docker-compose в фоновом режиме
local:
	docker-compose -f docker-compose.yml up -d

# Запуск приложения локально в dev-режиме (pprof на /debug/pprof)
run:
	go run cmd/gog/main.go

# Запуск всех тестов с подсчётом покрытия (в конце выводится общий % покрытия)
test:
	go test ./... -v -count=1 -coverprofile=coverage.out
	@echo "=== Покрытие тестами ==="
	@go tool cover -func=coverage.out | grep total | awk '{print "Общий процент покрытия: " $$3}'

# Генерация моков (требует mockgen в PATH).
generate:
	go generate ./internal/repository/...
