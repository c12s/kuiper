package startup

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/c12s/kuiper/repository"
	"github.com/c12s/kuiper/repository/etcd"
	"github.com/c12s/kuiper/service"
	"github.com/c12s/kuiper/startup/config"

	"github.com/c12s/kuiper/controller"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	Config *config.Config
	Logger *zap.Logger
}

func NewServer() *Server {
	logger, err := InitAndConfigureLogger()
	if err != nil {
		log.Printf("Error in initialization logger cause of: %s", err)
		panic(err)
	}
	server := &Server{
		Config: config.NewConfig(),
		Logger: logger,
	}

	server.Logger.Info("Server object successfully initialized", zap.Any("server", *server))
	return server
}

func (server Server) Start() {
	repository := server.initConfigRepository()
	service := server.initConfigService(repository)
	controller := server.initConfigController(*service)

	server.start(controller)
}

func (server *Server) initConfigRepository() repository.ConfigRepostory {

	repoLogger := server.Logger.Named("[ REPOSITORY ]")
	repo, err := etcd.New(repoLogger)
	if repo == nil || err != nil {
		repoLogger.Error("Repo object is null or error occured in initialization client , shutting down...",
			zap.Any("repoObject", repo),
			zap.Error(err),
		)
		os.Exit(1)
	}

	return repo
}

func (server *Server) initConfigService(repo repository.ConfigRepostory) *service.ConfigService {
	return service.New(repo, server.Logger.Named("[ SERVICE ]"))
}

func (server *Server) initConfigController(service service.ConfigService) *controller.Controller {
	return controller.New(service, server.Logger.Named("[ CONTROLLER ]"))
}

func (server Server) start(controller *controller.Controller) {
	router := mux.NewRouter()
	controller.Init(router)
	log := server.Logger.Named("[ SERVER ]")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", server.Config.SERVICE_PORT),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("error in serving http",
				zap.Error(err),
			)
		}
	}()

	gShutdownChannel := make(chan os.Signal, 1)
	signal.Notify(gShutdownChannel,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	<-gShutdownChannel

	timeout := time.Second * 15
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Error on shutting down server.")
	}

	log.Info("Gracefully shutdown executed.")
}
