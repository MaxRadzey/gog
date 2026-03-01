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

# Запуск всех тестов одной командой
test:
	go test ./... -v -count=1
