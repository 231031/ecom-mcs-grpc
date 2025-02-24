package main

import (
	"log"
	"net/http"

	"github.com/231031/ecom-mcs-grpc/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountUrl string `envconfig:"ACCOUNT_SERVICE_URL"`
	OrderUrl   string `envconfig:"ORDER_SERVICE_URL"`
	CatalogUrl string `envconfig:"CATALOG_SERVICE_URL"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	s, err := graphql.NewGraphQLServer(cfg.AccountUrl, cfg.CatalogUrl, cfg.OrderUrl)
	if err != nil {
		log.Fatal(err)
	}

	h := handler.NewDefaultServer(s.ToExecutablesSchema())
	p := playground.Handler("GraphQL", "/graphql")
	http.Handle("/playground", p)
	http.Handle("/graphql", h)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("graphql Listen and serve")
}
