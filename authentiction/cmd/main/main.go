package main

import (
	"log"
	"time"

	"github.com/231031/ecom-mcs-grpc/authentiction"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURl string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r authentiction.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = authentiction.NewRepository(cfg.DatabaseURl)
		if err != nil {
			log.Fatal(err)
		}
		return
	})
	defer r.Close()
	log.Println("Listening on port")

	s := authentiction.NewService(r)
	log.Fatal(authentiction.ListenGRPC(s, 50004))
}
