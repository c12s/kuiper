module github.com/c12s/kuiper

go 1.21.3

toolchain go1.21.4

require (
	github.com/c12s/magnetar v1.0.0
	github.com/c12s/oort v1.0.0
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/golang/protobuf v1.5.3
	github.com/nats-io/nats.go v1.31.0
	golang.org/x/exp v0.0.0-20230801115018-d63ba01acd4b
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.30.0
	iam-service v1.0.0
)

require (
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/nats-io/nkeys v0.4.5 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
)

replace github.com/c12s/magnetar => ../magnetar

replace iam-service => ../iam-service/iam-service

replace github.com/c12s/oort => ../oort
