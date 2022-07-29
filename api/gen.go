//nolint
//go:generate oapi-codegen -package api -generate types -o api.types.gen.go api.yaml
//go:generate oapi-codegen -package api -generate chi-server -o api.server.gen.go api.yaml
//go:generate oapi-codegen -package api -generate client -o api.client.gen.go api.yaml

package api
