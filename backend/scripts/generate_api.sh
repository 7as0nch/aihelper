#!/bin/bash

# 生成API代码
protoc --proto_path=. \
       --proto_path=./third_party \
       --go_out=paths=source_relative:. \
       --go-http_out=paths=source_relative:. \
       --go-grpc_out=paths=source_relative:. \
       api/aichat/v1/*.proto

echo "API代码生成完成"