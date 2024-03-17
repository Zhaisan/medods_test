up:
	docker-compose up -d

down:
	docker-compose down

clean:
	docker-compose down -v

logs:
	docker-compose logs -f mongodb

server:
	go run main.go

.PHONY: up down clean logs server
