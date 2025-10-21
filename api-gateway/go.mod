module api-gateway

go 1.23.0

require (
	github.com/go-chi/chi/v5 v5.2.3
	go-project/proto/auth v0.0.0
	go-project/proto/user v0.0.0
	google.golang.org/grpc v1.67.1
)

require (
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace (
	go-project/proto/auth => ../proto/auth
	go-project/proto/user => ../proto/user
)
