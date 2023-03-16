package server

import (
	"kuiper/service"

	"go.opentelemetry.io/otel/trace"
)

type ConfigHandler interface {
}

type configHandler struct {
	tracer        trace.Tracer
	configService service.ConfigService
}

func NewConfigHandler() configHandler {
	return configHandler{}
}
