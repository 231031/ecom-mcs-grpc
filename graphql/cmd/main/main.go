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
	AuthUrl       string `envconfig:"AUTH_SERVICE_URL"`
	AccountUrl    string `envconfig:"ACCOUNT_SERVICE_URL"`
	OrderUrl      string `envconfig:"ORDER_SERVICE_URL"`
	CatalogUrl    string `envconfig:"CATALOG_SERVICE_URL"`
	PublicKeyPath string `envconfig:"PUBLIC_KEY_PATH"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	middleware := graphql.NewAuthMiddleware(cfg.PublicKeyPath)
	s, err := graphql.NewGraphQLServer(cfg.AuthUrl, cfg.AccountUrl, cfg.CatalogUrl, cfg.OrderUrl)
	if err != nil {
		log.Fatal(err)
	}

	srv := s.ToExecutablesSchema(middleware)

	h := handler.NewDefaultServer(srv)
	p := playground.Handler("GraphQL", "/graphql")
	http.Handle("/playground", p)
	http.Handle("/graphql", graphql.ResponseWriterGetTokenMiddleware(h))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("graphql Listen and serve")
}
