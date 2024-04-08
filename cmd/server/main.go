package main

import (
	"log"
	"time"

	"github.com/c12s/kuiper/internal/domain"
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

	// client, errConn := clientv3.New(clientv3.Config{
	// 	Endpoints:   []string{"localhost:2379"},
	// 	DialTimeout: 5 * time.Second,
	// })
	// if errConn != nil {
	// 	log.Fatalln(errConn)
	// }
	// defer client.Close()

	// ctx := context.Background()

	// etcd := store.NewStandaloneConfigEtcdStore(client)

	paramSet := domain.NewParamSet("db_config", map[string]string{"port": "9999", "pass": "admin"})
	config := domain.NewStandaloneConfig("c12s", "v1.0.0", time.Now().Unix(), *paramSet)

	paramSet2 := domain.NewParamSet("db_config", map[string]string{"port": "1111"})
	config2 := domain.NewStandaloneConfig("c12s", "v1.0.0", time.Now().Unix(), *paramSet2)

	log.Println(config2.Diff(config))

	group1 := domain.NewConfigGroup("c12s", "g1", "v1.0.0", time.Now().Unix(), []domain.NamedParamSet{*paramSet2})
	group2 := domain.NewConfigGroup("c12s", "g1", "v1.0.0", time.Now().Unix(), []domain.NamedParamSet{*paramSet})

	log.Println(group2.Diff(group1))

	// err := etcd.Put(ctx, config)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
}
