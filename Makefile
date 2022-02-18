BINARY_NAME=gqsh
CMD_PREFIX=cmd/gqsh

build: ./target dep
	go build -o target/${BINARY_NAME} ${CMD_PREFIX}/main.go
	GOARCH=amd64 GOOS=darwin go build -o target/${BINARY_NAME}-darwin ${CMD_PREFIX}/main.go
	#GOARCH=amd64 GOOS=linux go build -o target/${BINARY_NAME}-linux ${CMD_PREFIX}/main.go
	#GOARCH=amd64 GOOS=window go build -o target/${BINARY_NAME}-windows ${CMD_PREFIX}/main.go

./target:
	@mkdir -p ./target

run: build
	./target/${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm -rf target

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod tidy

vet:
	go vet

lint: dev_deps
	golangci-lint run --enable-all

dev_deps:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
