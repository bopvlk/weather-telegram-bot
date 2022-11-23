.SILENT:

build:
	go mod download && go build -o bin/bot ./cmd/telegrambot/

docker-build: build
	docker build --rm -t telegram-bot .
	
run: docker-build
	docker run -d --name go-telegram-bot telegram-bot



lint:
	./scripts/linters.sh
.PHONY: lint

fix-imports:
	gogroup -order std,other,prefix=2.4-weather-forecast-bot --rewrite $(find . -type f -name "*.go" | grep -v /vendor/ |grep -v /.git/)

format:
	go vet ./...
	go fmt ./...