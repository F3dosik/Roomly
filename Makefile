.PHONY: up down down-clean migrate-down seed test
up:
	docker compose up --build -d

down:
	docker compose down

down-clean:
	docker compose down -v

migrate-down:
	docker compose --profile tools run --rm migrate-down

seed:
	docker compose --profile tools run --rm seed

test:
	go test ./...

test-e2e:
	go test -v -tags=e2e ./tests/e2e/...