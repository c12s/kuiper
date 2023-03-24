package server

import (
	"kuiper/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/trace"
)

type ConfigHandler interface {
	SaveConfig(c *gin.Context)
	GetConfig(c *gin.Context)
	DeleteConfig(c *gin.Context)
	CreateNewVersion(c *gin.Context)
	DeleteConfigsWithPrefix(c *gin.Context)
}

type configHandler struct {
	tracer        trace.Tracer
	logger        log.Logger
	configService service.ConfigService
	nats          nats.Conn
}

func NewConfigHandler(tracer trace.Tracer, logger log.Logger, configService service.ConfigService, conn nats.Conn) configHandler {
	return configHandler{tracer: tracer, logger: logger, configService: configService, nats: conn}
}
