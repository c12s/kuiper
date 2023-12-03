package startup

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"

	"github.com/c12s/kuiper/internal/configs"
	"github.com/c12s/kuiper/internal/servers"
	"github.com/c12s/kuiper/pkg/api"
	"github.com/c12s/kuiper/pkg/client/agent_queue"
	magnetarapi "github.com/c12s/magnetar/pkg/api"
	oortapi "github.com/c12s/oort/pkg/api"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type app struct {
	config                    *configs.Config
	grpcServer                *grpc.Server
	kuiperGrpcServer          api.KuiperServer
	evaluatorClient           oortapi.OortEvaluatorClient
	administratorClient       *oortapi.AdministrationAsyncClient
	magnetarClient            magnetarapi.MagnetarClient
	agentQueueClient          agent_queue.AgentQueueClient
	shutdownProcesses         []func()
	gracefulShutdownProcesses []func(wg *sync.WaitGroup)
}

func NewAppWithConfig(config *configs.Config) (*app, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	return &app{
		config:                    config,
		shutdownProcesses:         make([]func(), 0),
		gracefulShutdownProcesses: make([]func(wg *sync.WaitGroup), 0),
	}, nil
}

func (a *app) Start() error {
	a.init()
	return a.startGrpcServer()
}

func (a *app) GracefulStop(ctx context.Context) {
	// call all shutdown processes after a timeout or graceful shutdown processes completion
	defer a.shutdown()

	// wait for all graceful shutdown processes to complete
	wg := &sync.WaitGroup{}
	wg.Add(len(a.gracefulShutdownProcesses))

	for _, gracefulShutdownProcess := range a.gracefulShutdownProcesses {
		go gracefulShutdownProcess(wg)
	}

	// notify when graceful shutdown processes are done
	gracefulShutdownDone := make(chan struct{})
	go func() {
		wg.Wait()
		gracefulShutdownDone <- struct{}{}
	}()

	// wait for graceful shutdown processes to complete or for ctx timeout
	select {
	case <-ctx.Done():
		log.Println("ctx timeout ... shutting down")
	case <-gracefulShutdownDone:
		log.Println("app gracefully stopped")
	}
}

func (a *app) init() {
	natsConn, err := NewNatsConn(a.config.NatsAddress())
	if err != nil {
		log.Fatalln(err)
	}
	a.shutdownProcesses = append(a.shutdownProcesses, func() {
		log.Println("closing nats conn")
		natsConn.Close()
	})

	a.initMagnetarClient()
	a.initAgentQueueClient()
	a.initAdministratorClient()
	a.initEvaluatorClient()
	a.initKuiperGrpcServer(natsConn)
	a.initGrpcServer()
}

func (a *app) initGrpcServer() {
	if a.kuiperGrpcServer == nil {
		log.Fatalln("kuiper grpc server is nil")
	}
	s := grpc.NewServer()
	api.RegisterKuiperServer(s, a.kuiperGrpcServer)
	reflection.Register(s)
	a.grpcServer = s
}

func (a *app) initKuiperGrpcServer(conn *nats.Conn) {
	if a.magnetarClient == nil {
		log.Fatalln("magnetar client is nil")
	}
	if a.agentQueueClient == nil {
		log.Fatalln("blackhole client is nil")
	}
	if a.evaluatorClient == nil {
		log.Fatalln("evaluator client is nil")
	}
	if a.administratorClient == nil {
		log.Fatalln("administrator client is nil")
	}
	a.kuiperGrpcServer = servers.NewKuiperServer(conn, a.magnetarClient, a.evaluatorClient, a.administratorClient, a.agentQueueClient)
}

func (a *app) initEvaluatorClient() {
	client, err := newOortEvaluatorClient(a.config.OortAddress())
	if err != nil {
		log.Fatalln(err)
	}
	a.evaluatorClient = client
}

func (a *app) initAdministratorClient() {
	client, err := newOortAdministratorClient(a.config.NatsAddress())
	if err != nil {
		log.Fatalln(err)
	}
	a.administratorClient = client
}

func (a *app) initMagnetarClient() {
	client, err := newMagnetarClient(a.config.MagnetarAddress())
	if err != nil {
		log.Fatalln(err)
	}
	a.magnetarClient = client
}

func (a *app) initAgentQueueClient() {
	client, err := newAgentQueueClient(a.config.AgentQueueAddress())
	log.Printf("AgentQueue Address %s\n", a.config.AgentQueueAddress())
	log.Printf("%+v\n", client)
	if err != nil {
		log.Fatalln(err)
	}
	a.agentQueueClient = client
}

func (a *app) startGrpcServer() error {
	lis, err := net.Listen("tcp", a.config.ServerAddress())
	if err != nil {
		return err
	}
	go func() {
		log.Printf("server listening at %v", lis.Addr())
		if err := a.grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	a.gracefulShutdownProcesses = append(a.gracefulShutdownProcesses, func(wg *sync.WaitGroup) {
		a.grpcServer.GracefulStop()
		log.Println("grpc server gracefully stopped")
		wg.Done()
	})
	return nil
}

func (a *app) shutdown() {
	for _, shutdownProcess := range a.shutdownProcesses {
		shutdownProcess()
	}
}
