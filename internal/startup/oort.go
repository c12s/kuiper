package startup

import (
	oortapi "github.com/c12s/oort/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newOortEvaluatorClient(address string) (oortapi.OortEvaluatorClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return oortapi.NewOortEvaluatorClient(conn), nil
}

func newOortAdministratorClient(address string) (*oortapi.AdministrationAsyncClient, error) {
	return oortapi.NewAdministrationAsyncClient(address)
}
