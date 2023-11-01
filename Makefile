proto:
	rm -f Bankingpb/*go
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	Bankingpb/Golang-gRPC-Banking.proto

proto2:
	rm -f Bankingpb/*go
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=.\
    --grpc-gateway_opt=paths=source_relative \
	Bankingpb/Golang-gRPC-Banking.proto

gateway:
	
	protoc -I . --grpc-gateway_out . \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt grpc_api_configuration=path/to/config.yaml \
		--grpc-gateway_opt standalone=true \
		Bankingpb/Golang-gRPC-Banking.proto


server:
	go run Banking_server/main.go

client:
	go run Banking_client/main.go

proxy:
	go run Banking_proxy/proxy.go