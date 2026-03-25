.PHONY: up down down-clean migrate-down seed test test-e2e
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
	docker compose run --rm tests

test-e2e:
	docker compose run --rm tests-e2e 