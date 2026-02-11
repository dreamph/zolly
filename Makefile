.PHONY: proto build build-plugin-auth build-plugin-hmac build-plugins

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		plugin/proto/middleware.proto

build:
	go build -o zolly .

build-plugin-auth:
	cd examples/plugins/auth && go build -o ../../../plugins/auth-plugin .

build-plugin-hmac:
	cd examples/plugins/hmac && go build -o ../../../plugins/hmac-plugin .

build-plugins: build-plugin-auth build-plugin-hmac
