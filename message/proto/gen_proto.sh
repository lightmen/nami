#!/bin/bash

protoc --go-grpc_opt=require_unimplemented_servers=false --gogoslick_out=paths=source_relative:.. --go-grpc_out=paths=source_relative:..  --proto_path=. *.proto