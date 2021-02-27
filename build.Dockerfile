# FROM golang:1.16-alpine AS builder
FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    GOPROXY=https://proxy.golang.org \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPRIVATE=

RUN apk --update upgrade \
    && apk --no-cache --no-progress add git bash musl-dev curl tar ca-certificates tzdata \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*
    
ENV APPPATH="./cmd/app"    

# Move to working directory /build
WORKDIR /build

# RUN export VERSION="$(shell git describe --abbrev=0)"

# RUN go get github.com/go-bindata/go-bindata/...
# RUN chmod +x ${GOPATH}/bin/go-bindata
RUN go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen
RUN chmod +x ${GOPATH}/bin/oapi-codegen

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Generate files
# RUN go generate ./...

# Download misspell binary to bin folder in $GOPATH
RUN  curl -sfL https://raw.githubusercontent.com/client9/misspell/master/install-misspell.sh | bash -s -- -b $(go env GOPATH)/bin v0.3.4

# Download golangci-lint binary to bin folder in $GOPATH
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.37.0
RUN golangci-lint run --out-format=tab --tests=false ./...

RUN go test ./... -tags=skip

# Build the application
RUN go build -ldflags "-X main.version=develop" -o app ${APPPATH}

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/app .

# Build a small image
FROM alpine:latest 

RUN apk --no-cache add ca-certificates

COPY --from=builder /dist/app /

# Open port
EXPOSE 8080

# Command to run
ENTRYPOINT ["/app"]