# Go language blank application

Some Go language blank application.

## Launch generation

```sh
go generate ./...
```

or 

```sh
make gen
```

## Launch  tests

```sh
go test ./...
```

or 

```sh
make test
```

## Docker-compose

```sh
docker-compose up -d someconteiner
docker-compose up --build app
```

or 

```sh
# Start all containers.
make up

# Start some containers.
make up-some

# Start APP container.
make up-app

# Stop containers.
make down
```

## Launch application

```sh
APP_DB_HOST=127.0.0.1
APP_DB_PORT=5432
APP_DB_NAME=database
APP_DB_USER=$DB_USER
APP_DB_PASSWORD=$DB_PASSWORD
APP_HOST=
APP_PORT=80
APP_PATH=/

go run ./cmd/app/*
```

or 

```sh
# Launch application.
make run

# Show help.
make help
```

