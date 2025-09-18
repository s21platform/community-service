env:
	go run cmd/system/autoganarate_env/main.go

protogen:
	protoc --go_out=. --go-grpc_out=. ./api/community.proto
	protoc --doc_out=. --doc_opt=markdown,GRPC_API.md ./api/community.proto

codegen:
	oapi-codegen -generate chi-server -package apigen api/schema.yaml > internal/generated/server.gen.go
	oapi-codegen -generate types -package apigen api/schema.yaml > internal/generated/models.gen.go