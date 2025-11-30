package graphql

import (
	"github.com/231031/ecom-mcs-grpc/account"
	"github.com/231031/ecom-mcs-grpc/authentication"
	"github.com/231031/ecom-mcs-grpc/catalog"
	"github.com/231031/ecom-mcs-grpc/order"
	"github.com/99designs/gqlgen/graphql"
)

//go:generate go run github.com/99designs/gqlgen generate

type Server struct {
	authClient    *authentication.Client
	accountClient *account.Client
	catalogClient *catalog.Client
	orderClient   *order.Client
}

func NewGraphQLServer(authUrl, accountUrl, catalogUrl, orderUrl string) (*Server, error) {
	authClient, err := authentication.NewClient(authUrl)
	if err != nil {
		authClient.Close()
		return nil, err
	}

	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		authClient.Close()
		accountClient.Close()
		return nil, err
	}

	catalogClient, err := catalog.NewClient(catalogUrl)
	if err != nil {
		accountClient.Close()
		return nil, err
	}

	orderClient, err := order.NewClient(orderUrl)
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return nil, err
	}

	return &Server{
		authClient,
		accountClient,
		catalogClient,
		orderClient,
	}, nil
}

func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}

}

func (s *Server) ToExecutablesSchema(m *authMiddlewre) graphql.ExecutableSchema {
	c := Config{Resolvers: s}
	c.Directives.HasRole = m.HasRole

	return NewExecutableSchema(c)
}
