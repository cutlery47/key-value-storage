.PHONY: client storage all

client:
	mkdir client/build
	go build -C client/cmd -o ../build/app main.go

storage:
	mkdir storage/build
	go build -C storage/cmd -o ../build/app main.go

all:
	mkdir client/build
	mkdir storage/build
	go build -C client/cmd -o ../build/app main.go
	go build -C storage/cmd -o ../build/app main.go