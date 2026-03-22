up:
	docker compose up --build -d

down:
	docker compose down

seed:
	docker compose exec roomly ./roomly seed

test:
	go test ./...