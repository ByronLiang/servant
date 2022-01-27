cd public/protocol
protoc --go_out=plugins=grpc,paths=source_relative:../pb *.proto
