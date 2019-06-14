all: build

build:
	go build -o bin/yuancast ./cmd/go-mysql-elasticsearch
	bin/yuancast
