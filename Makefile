.PHONY: 

APPLICATION = "app"

SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse HEAD)
VERSION_GIT := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
VERSION := $(if $(VERSION),$(VERSION),$(VERSION_GIT))
GIT_BRANCH := $(subst heads/,,$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null))

mod:
	go mod tidy

gen:
	go generate ./...

test: mod gen
	go test ./... -v -count=1 -race	

## Format the Code.
fmt:
	gofmt -s -l -w $(SRCS)

run: mod gen
	export $(cut -d= -f1 .env) && go run -ldflags "-X main.version=$(VERSION)" ./cmd/$(APPLICATION)/...

build: mod gen
	go build -ldflags "-X main.version=$(VERSION)" -o dist/$(APPLICATION) ./cmd/$(APPLICATION)/...

help: mod gen
	go run -ldflags "-X main.version=$(VERSION) main.commit=$(SHA)" ./cmd/$(APPLICATION)/* run --help

# Docker-compose operations.

# Up server environment.
up-env:
	# docker-compose pull
	docker-compose up -d db

# Up server.
up-app:
	docker-compose up --build app

# Up server and environment.
up: up-env up-app

down:
	docker-compose down -v --remove-orphans

# Drop server data.
drop: down
	# rmdir -rf ./.tmp/pgdata

# Drop and up server.
drop-up: drop up

# Drop and up server environment.
drop-up-env: drop up-env

