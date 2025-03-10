package graphql

import (
	"context"
	"errors"
	"log"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Buyer(ctx context.Context, id string) (*AccountBuyer, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != "" {
		_, err := r.server.accountClient.GetAccount(ctx, id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return nil, nil
	}

	// skip := uint64(0)
	// take := uint64(0)
	// accountLists, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }

	var accounts *AccountBuyer
	return accounts, nil
}

func (r *queryResolver) Sellers(ctx context.Context, pagination *PaginationInput, id *string) ([]*AccountSeller, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		_, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*AccountSeller{}, nil
	}

	// skip := uint64(0)
	// take := uint64(0)
	// accountLists, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }

	var accounts []*AccountSeller

	return accounts, nil
}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	q := ""
	if query != nil {
		q = *query
	}

	if id != nil && q == "" {
		p, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*Product{{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}
	productList, err := r.server.catalogClient.GetProducts(ctx, skip, take, nil, q)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []*Product
	for _, p := range productList {
		products = append(products,
			&Product{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			},
		)
	}
	return products, nil
}

func (r *queryResolver) Orders(ctx context.Context, id *string) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id == nil {
		return nil, errors.New("id is required")
	}

	orderList, err := r.server.orderClient.GetOrdersForAccount(ctx, *id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	orders := []*Order{}
	for _, o := range orderList {
		var products []*OrderProduct
		for _, p := range o.Products {
			products = append(products, &OrderProduct{
				Product: &Product{
					ID:          p.ID,
					Name:        p.Name,
					Description: p.Description,
					Price:       p.Price,
				},
				Quantity: int(p.Quantity),
			})
		}

		orders = append(orders, &Order{
			ID:         o.ID,
			TotalPrice: o.TotalPrice,
			CreatedAt:  o.CreatedAt,
			Products:   products,
		})
	}
	return orders, nil
}

func (p *PaginationInput) bounds() (uint64, uint64) {
	skipVal := uint64(0)
	takeVal := uint64(0)

	if p != nil {
		skipVal = uint64((*p).Skip)
		takeVal = uint64((*p).Take)
	}

	return skipVal, takeVal
}
