module user-service

go 1.23.0

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	go-project/proto/user v0.0.0
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.36.6
)

require (
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
)

replace go-project/proto/user => ../proto/user
