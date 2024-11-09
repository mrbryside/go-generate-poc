prepare:
	brew install make
	go install github.com/swaggo/swag/cmd/swag@latest
	go install golang.org/x/tools/cmd/goimports@latest

run:
	go run test/main.go

gen:
	go run main.go generate test

swag:
	cd test && swag init
