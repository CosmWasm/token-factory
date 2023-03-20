#!/usr/bin/env bash

#== Requirements ==
#
## make sure your `go env GOPATH` is in the `$PATH`
## Install:
## + latest buf (v1.0.0-rc11 or later)
## + protobuf v3
#
## All protoc dependencies must be installed not in the module scope
## currently we must use grpc-gateway v1
# cd ~
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0
# go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest
# go get github.com/regen-network/cosmos-proto@latest # doesn't work in install mode
# go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@v0.3.1

# To run this, i had to modify the following:
# git clone https://github.com/regen-network/cosmos-proto.git
#
# Modify interfacetype/interfacetype.go and comment out the following lines:
# 
# if len(message.OneofDecl) != 1 {
# 	panic("interfacetype only supports messages with exactly one oneof declaration")
# }
# for _, field := range message.Field {
# 	if idx := field.OneofIndex; idx == nil || *idx != 0 {
# 		panic("all fields in interfacetype message must belong to the oneof")
# 	}
# }
#
# then:
# cd cosmos-proto/protoc-gen-gocosmos
# go install .

set -eo pipefail

echo "Generating gogo proto code"
cd proto
buf mod update
cd ..
buf generate

# move proto files to the right places
cp -r ./github.com/CosmWasm/token-factory/x/* x/
rm -rf ./github.com

go mod tidy 


