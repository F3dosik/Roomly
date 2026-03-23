up:
	docker compose up --build -d

down:
	docker compose down

migrate-down:
	docker compose --profile tools run --rm migrate-down

seed:
	docker compose exec roomly ./roomly seed

test:
	go test ./...