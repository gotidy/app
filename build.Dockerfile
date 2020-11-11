FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    GOPROXY=https://proxy.golang.org \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ENV APPPATH="./cmd/app"    

# Move to working directory /build
WORKDIR /build


RUN go get github.com/go-bindata/go-bindata/...
RUN chmod +x ${GOPATH}/bin/go-bindata
RUN go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen
RUN chmod +x ${GOPATH}/bin/oapi-codegen

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Generate files
RUN go generate 

RUN go test ./...
# Load linter
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.27.0
RUN golangci-lint run --out-format=tab --tests=false ./...

# Build the application
RUN go build -ldflags "-X main.version=develop" -o app ./cmd/app

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