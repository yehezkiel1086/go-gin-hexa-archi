composeup:
	docker-compose up -d

composedown:
	docker-compose down

run:
	go run cmd/main.go

.PHONY: composeup composedown run
