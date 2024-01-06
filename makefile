PROTO_SRC_FILES_DATA=$(shell find ./proto/data -type f -name "*.proto" | sed 's/\/proto//g')
PROTO_SRC_FILES_SERVICE=$(shell find ./proto/service -type f -name "*.proto" | sed 's/\/proto//g')

.PHONY: proto
proto:
	rm ./data/*.pb.go; \
	rm ./service/*.pb.go; \
	cd ./proto; \
	protoc -I=. \
		--go_out=paths=source_relative:../  \
		$(PROTO_SRC_FILES_DATA); \
	protoc -I=. \
		-I=${GOPATH}/src/github.com/googleapis \
		-I=${GOPATH}/src/github.com/grpc-gateway \
		--go_out=paths=source_relative:../  \
		--go-grpc_out=../ --go-grpc_opt=paths=source_relative \
		--go-grpc_opt=require_unimplemented_servers=false \
		--grpc-gateway_out=../ \
		--grpc-gateway_opt=logtostderr=true \
		--grpc-gateway_opt=paths=source_relative \
    --grpc-gateway_opt=generate_unbound_methods=true \
    --openapiv2_opt=logtostderr=true \
    --openapiv2_out=allow_merge=true,merge_file_name=../doc/spec/apidocs:. \
		$(PROTO_SRC_FILES_SERVICE)

install:
	cd ./cmd/signcli/; \
	go install -tags 'main' -mod=readonly

.PHONY: test
test:
	MallocNanoZone=0 go test -race -timeout 60s ./...

bench:
	go test ./... -bench=. -benchtime=10s

fmt:
	go fmt ./...
	clang-format -i ./proto/**/*.proto

lint:
	go vet ./...

build:
	cd ./cmd/signcli/; \
	go build -o signcli -gcflags '-m' -tags 'main'

.PHONY: doc
doc:
	docker-compose up -d docs

docdown:
	docker-compose down

init:
	go run cmd/signcli/* init --home ./.signcli

run:
	go run cmd/signcli/* serve --home ./.signcli --gateway-endpoint localhost:8080 --log-level info
