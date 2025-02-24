package account

import (
	"context"
	"log"

	"github.com/231031/ecom-mcs-grpc/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := pb.NewAccountServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.service.PostAccount(
		ctx,
		&pb.PostAccountRequest{Name: name},
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Account{
		Name: r.Account.Name,
		ID:   r.Account.Id,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	r, err := c.service.GetAccount(
		ctx,
		&pb.GetAccountRequest{Id: id},
	)
	if err != nil {
		return nil, err
	}

	return &Account{
		Name: r.Account.Name,
		ID:   r.Account.Id,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	r, err := c.service.GetAccounts(
		ctx,
		&pb.GetAccountsRequest{Skip: skip, Take: take},
	)
	if err != nil {
		return nil, err
	}

	accounts := []Account{}
	for _, a := range r.Accounts {
		accounts = append(accounts, Account{
			Name: a.Name,
			ID:   a.Id,
		})
	}

	return accounts, nil
}
