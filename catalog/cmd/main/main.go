package main

import (
	"log"
	"time"

	"github.com/231031/ecom-mcs-grpc/catalog"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURl     string `envconfig:"DATABASE_URL"`
	ElasticUsername string `envconfig:"ELASTIC_USERNAME"`
	ElasticPassword string `envconfig:"ELASTIC_PASSWORD"`
}

func main() {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(cfg.DatabaseURl, cfg.ElasticUsername, cfg.ElasticPassword)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	})
	defer r.Close()
	log.Println("Listening on port")

	s := catalog.NewService(r)
	log.Fatal(catalog.ListenGRPC(s, 50002))
}
