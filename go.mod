module github.com/c12s/kuiper

go 1.21

require (
	github.com/c12s/magnetar v1.0.0
	github.com/c12s/oort v1.0.0
	github.com/golang/protobuf v1.5.3
	github.com/nats-io/nats.go v1.28.0
	golang.org/x/exp v0.0.0-20230801115018-d63ba01acd4b
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/nats-io/nats-server/v2 v2.9.21 // indirect
	github.com/nats-io/nkeys v0.4.4 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
)

replace github.com/c12s/magnetar => ../magnetar

replace github.com/c12s/oort => ../oort
