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

func (c *Client) PostAccountBuyer(ctx context.Context, in *pb.PostAccountBuyerRequest) (*Buyer, error) {
	r, err := c.service.PostAccountBuyer(ctx, in)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Buyer{
		ID: r.Account.Id,
		BaseInfo: BaseInfo{
			FirstName: r.Account.BaseInfo.FirstName,
			LastName:  r.Account.BaseInfo.LastName,
			Email:     r.Account.BaseInfo.Email,
			Phone:     r.Account.BaseInfo.Phone,
			Address:   r.Account.BaseInfo.Address,
		},
	}, nil
}

func (c *Client) PostAccountSeller(ctx context.Context, in *pb.PostAccountSellerRequest) (*Seller, error) {
	r, err := c.service.PostAccountSeller(ctx, in)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Seller{
		ID:        r.Account.Id,
		StoreName: r.Account.StoreName,
		BaseInfo: BaseInfo{
			FirstName: r.Account.BaseInfo.FirstName,
			LastName:  r.Account.BaseInfo.LastName,
			Email:     r.Account.BaseInfo.Email,
			Phone:     r.Account.BaseInfo.Phone,
			Address:   r.Account.BaseInfo.Address,
		},
	}, nil
}

func (c *Client) UpdateAccountBuyer(ctx context.Context, in *pb.AccountBuyer) (*pb.AccountBuyer, error) {
	_, err := c.service.UpdateAccountBuyer(ctx, in)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func (c *Client) UpdateAccountSeller(ctx context.Context, in *pb.AccountSeller) (*pb.AccountSeller, error) {
	_, err := c.service.UpdateAccountSeller(ctx, in)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	r, err := c.service.GetAccountBuyer(
		ctx,
		&pb.GetAccountBuyerRequest{Id: id},
	)
	if err != nil {
		return nil, err
	}

	return &Account{
		ID: r.Id,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	r, err := c.service.GetAccountsSeller(
		ctx,
		&pb.GetAccountsSellerRequest{Skip: skip, Take: take},
	)
	if err != nil {
		return nil, err
	}

	accounts := []Account{}
	for _, a := range r.Accounts {
		accounts = append(accounts, Account{
			ID: a.Id,
		})
	}

	return accounts, nil
}
