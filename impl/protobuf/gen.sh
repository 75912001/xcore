#!/bin/bash

genMod=go_out      #google.golang.org/protobuf/cmd/protoc-gen-go
#genMod=gofast_out  #github.com/gogo/protobuf/protoc-gen-gofast

########################################################################################################################
#gateway
srcFile="./proto/error.code.proto"
protoc --${genMod}=. ${srcFile}

echo "生成 protobuf 完成"
#read input

exit 0
