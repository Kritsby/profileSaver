BINARY_NAME=run_service

build:
	go build -o ${BINARY_NAME} cmd/service.go

run:
	go build -o ${BINARY_NAME} cmd/service.go
	./${BINARY_NAME}

tests:
	go test	./...

clean:
	go clean
	rm ${BINARY_NAME}