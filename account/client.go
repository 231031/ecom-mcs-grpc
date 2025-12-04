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

func NewClient(url string, opts ...grpc.DialOption) (*Client, error) {
	defaultOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	finalOpts := append(defaultOpts, opts...)

	conn, err := grpc.NewClient(url, finalOpts...)
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

func (c *Client) GetAccountBuyerByID(ctx context.Context, id string) (*Buyer, error) {
	r, err := c.service.GetAccountBuyer(
		ctx,
		&pb.GetAccountRequest{Id: id},
	)
	if err != nil {
		return nil, err
	}

	return &Buyer{
		ID: r.Id,
		BaseInfo: BaseInfo{
			FirstName: r.BaseInfo.FirstName,
			LastName:  r.BaseInfo.LastName,
			Phone:     r.BaseInfo.Phone,
			Address:   r.BaseInfo.Address,
		},
	}, nil
}

func (c *Client) GetAccountSellerByID(ctx context.Context, id string) (*Seller, error) {
	r, err := c.service.GetAccountSeller(
		ctx,
		&pb.GetAccountRequest{Id: id},
	)
	if err != nil {
		return nil, err
	}

	return &Seller{
		ID:        r.Id,
		StoreName: r.StoreName,
		BaseInfo: BaseInfo{
			FirstName: r.BaseInfo.FirstName,
			LastName:  r.BaseInfo.LastName,
			Phone:     r.BaseInfo.Phone,
			Address:   r.BaseInfo.Address,
		},
	}, nil
}

func (c *Client) GetAccountSellers(ctx context.Context, ids []string, skip uint64, take uint64) ([]Seller, error) {
	r, err := c.service.GetAccountSellers(
		ctx,
		&pb.GetAccountSellersRequest{Ids: ids, Skip: skip, Take: take},
	)
	if err != nil {
		return nil, err
	}

	sellers := []Seller{}
	for _, a := range r.Accounts {
		sellers = append(sellers, Seller{
			ID:        a.Id,
			StoreName: a.StoreName,
			BaseInfo: BaseInfo{
				FirstName: a.BaseInfo.FirstName,
				LastName:  a.BaseInfo.LastName,
				Phone:     a.BaseInfo.Phone,
				Address:   a.BaseInfo.Address,
			},
		})
	}

	return sellers, nil
}
