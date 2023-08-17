GO_WORKSPACE := $(GOPATH/src)

protoc:
	protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=. --experimental_allow_proto3_optional
	@echo "Protoc success compiled"

# protoc:
# 	protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.
# 	@echo "Protoc success compiled"

# protoc:
# 	protoc --proto_path=proto proto/*.proto --go_out=$(GO_WORKSPACE) --go-grpc_out=$(GO_WORKSPACE)
# 	@echo "Protoc success compiled"