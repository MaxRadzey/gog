# Генерация Swagger-доки
swag:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -d . --parseInternal -o docs

# Сборка CLI-клиента под текущую ОС
build-client:
	go build -o bin/gog-client ./cmd/client

# Сборка клиента под Linux (amd64)
build-client-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/gog-client-linux ./cmd/client

# Сборка клиента под Windows (amd64)
build-client-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/gog-client-windows.exe ./cmd/client

# Сборка клиента под macOS (Intel и Apple Silicon)
build-client-apple:
	GOOS=darwin GOARCH=arm64 go build -o bin/gog-client-apple ./cmd/client

# Форматирование кода + сортировка импортов
fmt:
	go run golang.org/x/tools/cmd/goimports@latest -w .

# Запуск приложения локально через docker-compose в фоновом режиме
local:
	docker-compose -f docker-compose.yml up -d

# Запуск приложения локально в dev-режиме (pprof на /debug/pprof)
run:
	go run ./cmd/server

# Исключения из покрытия берём из codecov.yml (ignore)
COVER_EXCLUDE := $(shell grep -E '^\s+-\s+"' codecov.yml 2>/dev/null | sed 's/.*"\([^"]*\)".*/\1/' | tr -d '*' | tr '/' '\n' | grep -v '^$$' | sort -u | sed 's/^/\//' | sed 's/$$/|/' | tr -d '\n' | sed 's/|$$//')
COVER_PKGS := $(shell go list ./... | grep -E -v '$(COVER_EXCLUDE)' | tr '\n' ',' | sed 's/,$$//')

# Запуск всех тестов с подсчётом покрытия (в конце выводится общий % покрытия)
test:
	go test ./... -v -count=1 -coverprofile=coverage.out -coverpkg=$(COVER_PKGS)
	@echo ""
	@echo "=== Покрытие тестами ==="
	@go tool cover -func=coverage.out | grep total | awk '{print "Общий процент покрытия: " $$3}'

# Генерация моков (требует mockgen в PATH).
generate:
	go generate ./internal/repository/...
