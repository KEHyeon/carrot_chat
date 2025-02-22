.PHONY: proto
PROTO_SRC = proto/chat.proto
PROTO_OUT = pkg/chat/pb
proto:
	@echo "🔄 Compiling Protobuf files..."
	protoc --go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
	       --go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
	       $(PROTO_SRC)