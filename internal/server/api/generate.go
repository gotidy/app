//go:generate echo "oapi-codegen"
//go:generate oapi-codegen -generate types -package api -o types.gen.go openapi.yaml
//go:generate oapi-codegen -generate chi-server,spec -package api -o server.gen.go openapi.yaml
//go:generate goimports -w ./types.gen.go
//go:generate goimports -w ./server.gen.go
package api
