run:
	go run ./cmd/main.go

build:
	go build ./cmd/main.go

compose-up:
	docker-compose build
	docker-compose up -d

compose-down:
	docker-compose down

restart-app:
	docker-compose restart verve-app

load-test:
	k6 run ./load-test-script.js