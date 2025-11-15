package main

import (
	"log"

	"github.com/231031/ecom-mcs-grpc/authentication"
	"github.com/231031/ecom-mcs-grpc/authentication/model"
	"github.com/231031/ecom-mcs-grpc/authentication/repository"
	"github.com/231031/ecom-mcs-grpc/authentication/service"
	"github.com/231031/ecom-mcs-grpc/authentication/utils"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	var cfg model.Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	db, err := authentication.ConnectPostgres(cfg.DatabaseURl)
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := authentication.ConnectRedis(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatal(err)
	}

	r := repository.NewRepository(db, redisClient)
	defer r.Close()
	log.Println("Listening on port")

	tokenCfg := utils.ConfigGenerateKey(&cfg)
	tokenService := service.NewTokenService(r, tokenCfg)
	s := service.NewService(r, tokenService)
	log.Fatal(authentication.ListenGRPC(s, 50004))
}
