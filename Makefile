
test:
	docker-compose --file e2e-docker-compose.yml up --build --abort-on-container-exit

run:
	docker-compose -f contained-docker-compose.yml up --build

build:
	docker build --tag redis-proxy .