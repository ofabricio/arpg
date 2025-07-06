
run:
	go run cmd/main.go

build:
	go build -ldflags "-H=windowsgui" cmd/main.go
