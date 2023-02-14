protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protos/payload.proto
go build -o bin/payload ./server/ctl
go build -o bin/persist ./server/persist/ctl/main.go
go build -o bin/client ./terminal_client/main.go
