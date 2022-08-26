.PHONY:

build:
	go build -o ./.bin/main cmd/main/app.go
run: build
	./.bin/main