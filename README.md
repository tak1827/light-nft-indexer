# light-nft-indexer
The EVM compatible light weight NFT Indexer, no DB and utilize free account of infra


## Install Proto Toools
Install protocol buffer compiler. Ref: https://grpc.io/docs/protoc-installation/
```sh
# Linux
apt install -y protobuf-compiler
protoc --version  # Ensure compiler version is 3+

# Mac
brew install protobuf
protoc --version  # Ensure compiler version is 3+
```
Install grpc. Ref: https://grpc.io/docs/languages/go/quickstart/
```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```
Install grpc gateway. Ref: https://github.com/grpc-ecosystem/grpc-gateway
```sh
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.0
```
Clone repos.
```sh
cd ${GOPATH}/src/github.com/
# clone googleapis
git clone https://github.com/googleapis/googleapis.git
# clone grpc gateway
git clone --branch v2.19.0 https://github.com/grpc-ecosystem/grpc-gateway.git
```