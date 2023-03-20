package main

import (
	"context"
	"kuiper/server"
	"kuiper/service"
	"kuiper/store"
	"kuiper/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	logger := log.Default()

	ctx := context.Background()
	//init exporter
	exporter, err := util.NewJaegerExporter("http://127.0.0.1:14268/api/traces")
	if err != nil {
		logger.Fatalf(err.Error())
	}
	//init traceprovider
	tp := util.NewTraceProvider(exporter)
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)
	tracer := tp.Tracer("kuiper")
	otel.SetTextMapPropagator(propagation.TraceContext{})

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware("kuiper"))

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 10 * time.Second,
	})

	cfgStore := store.NewConfigStore(*cli, *logger, tracer)
	cfgService := service.NewConfigService(cfgStore, *logger, tracer)
	handler := server.NewConfigHandler(tracer, *logger, cfgService)

	router.POST("/api/config", handler.SaveConfig)
	router.GET("/api/config/:id/:ver", handler.GetConfig)
	router.POST("/api/config/:id/", handler.CreateNewVersion)

	// start server
	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
