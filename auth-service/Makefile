generate:
	make generate-grpc

generate-grpc:
	protoc --go_out=pkg --go_opt=paths=source_relative --go-grpc_out=pkg \
	 --go-grpc_opt=paths=source_relative api/auth_v1/auth.proto