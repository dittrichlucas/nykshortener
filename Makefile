up:
	docker-compose up -d

down:
	docker-compose down

init:
	make down && make up

