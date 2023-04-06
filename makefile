BINARY_NAME=run_service

build:
	go build -o ${BINARY_NAME} cmd/service.go

run:
	go run cmd/service.go

tests:
	go test	./...

clean:
	go clean
	rm ${BINARY_NAME}