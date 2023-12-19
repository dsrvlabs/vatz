.PHONY: test build_proto build_proto_gw

test:
	go test -v

build_proto:
	@protoc -I=. \
		--go_out ./ --go_opt paths=import \
		--go-grpc_out ./ --go-grpc_opt paths=import \
		proto/*.proto

build_proto_gw:
	@protoc -I=. \
		--grpc-gateway_out ./ \
		proto/*.proto
