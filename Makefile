# Commands to run inside the container:
build-api:
	CGO_ENABLED=0 GOOS=linux go build -o littlejohn ./cmd

build-tests:
	CGO_ENABLED=0 GOOS=linux go test -c ./tests -o integration_tests

test-api:
	go test ./internal/...

# Commands to run outside the container:
build:
	docker build -t littlejohn .

run:
	docker run --name littlejohn_local --rm -d -p "8080:8080" -e "PORT=8080" littlejohn
	docker ps

integration-tests:
	docker exec -it littlejohn_local /app/tests/integration_tests

stop:
	docker stop littlejohn_local
