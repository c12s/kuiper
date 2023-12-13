package startup

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	apolloapi "iam-service/proto1"
)

func newApolloClient(address string) (apolloapi.AuthServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return apolloapi.NewAuthServiceClient(conn), nil
}
