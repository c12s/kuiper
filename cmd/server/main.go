package main

import (
	"context"
	"log"
	"time"

	"github.com/c12s/kuiper/internal/domain"
	"github.com/c12s/kuiper/internal/store"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	// config, err := configs.NewFromEnv()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// app, err := startup.NewAppWithConfig(config)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// err = app.Start()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// shutdown := make(chan os.Signal, 1)
	// signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	// <-shutdown

	// timeout := 10 * time.Second
	// ctx, cancel := context.WithTimeout(context.Background(), timeout)
	// defer cancel()
	// app.GracefulStop(ctx)

	client, errConn := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if errConn != nil {
		log.Fatalln(errConn)
	}
	defer client.Close()

	ctx := context.Background()

	etcd := store.NewStandaloneConfigEtcdStore(client)

	paramSet := domain.NewParamSet("db_config2", map[string]string{"port": "9999", "pass": "admin"})
	config := domain.NewStandaloneConfig("c12s", "v1.0.0", time.Now().Unix(), *paramSet)

	err := etcd.Put(ctx, config)
	if err != nil {
		log.Fatalln(err)
	}
}
