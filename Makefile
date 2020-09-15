.EXPORT_ALL_VARIABLES:

COMPOSE_CONVERT_WINDOWS_PATHS=1

tidy:
	go mod tidy
fmt:
	go fmt ./...
prepare_lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0
lint: prepare_lint fmt tidy
	golangci-lint run ./...
run:
	go run main.go --debug run
build:
	go build -o fanouter.exe main.go
test:
	go test -race ./...
testv:
	go test -v -race ./...
testi:
	go test -v -race ./tests -tags=integration