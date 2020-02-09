dev:
	go run main.go
build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .
docker:
	docker build -t maps-api-go .
