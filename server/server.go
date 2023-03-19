package server

import (
	"kuiper/service"
	"log"

	"go.opentelemetry.io/otel/trace"
)

type ConfigHandler interface {
}

type configHandler struct {
	tracer        trace.Tracer
	logger        log.Logger
	configService service.ConfigService
}

func NewConfigHandler(tracer trace.Tracer, logger log.Logger, configService service.ConfigService) configHandler {
	return configHandler{tracer: tracer, logger: logger, configService: configService}
}
