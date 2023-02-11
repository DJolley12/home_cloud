protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protos/payload.proto
go build ./server/payload/ctl
go build ./terminal_client/main.go
