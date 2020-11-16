VERSION := "$(shell git describe --abbrev=0)"
MODULE := "$(shell git config --get remote.origin.url | sed 's|^https\://\([^ <]*\)\(.*\)\.git|\1|g')"
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(dir $(MKFILE_PATH))
update-pkg-cache:
	cd $$HOME && \
	GOPROXY=https://proxy.golang.org GO111MODULE=on go get $(MODULE)@$(VERSION)

gen:
	go generate ./...

test:
	go test ./...	

up:
	docker-compose up -d somecontainer

	docker-compose up --build app	

up-devbox:
	docker-compose up -d somecontainer

up-app:
	docker-compose up --build app

down:
	docker-compose down -v --remove-orphans

run:
	source .env && go run ./cmd/app/*

help:
	go run ./cmd/app/* run --help