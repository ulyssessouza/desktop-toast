run:
	go run main.go

build:
	go build -o ./bin/toast-test .

test:
	gotestsum ./toast/tests/...
