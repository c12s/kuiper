package server

import (
	"kuiper/service"
	"log"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type ConfigHandler interface {
	SaveConfig(c *gin.Context)
	GetConfig(c *gin.Context)
}

type configHandler struct {
	tracer        trace.Tracer
	logger        log.Logger
	configService service.ConfigService
}

func NewConfigHandler(tracer trace.Tracer, logger log.Logger, configService service.ConfigService) configHandler {
	return configHandler{tracer: tracer, logger: logger, configService: configService}
}
