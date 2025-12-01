package graphql

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/231031/ecom-mcs-grpc/account/pb"
	auth_pb "github.com/231031/ecom-mcs-grpc/authentication/pb"
	"github.com/231031/ecom-mcs-grpc/order"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrInvalidInfo      = errors.New("invalid info")
)

type mutationResolver struct {
	server *Server
}

func (m *mutationResolver) CreateUser(ctx context.Context, email string, password string, role RoleType) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	in := &auth_pb.CreateUserRequest{
		Email:    email,
		Password: password,
		Role:     MapRoleToInt(role),
	}

	_, err := m.server.authClient.CreateUser(ctx, in)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return "", nil
}

func (m *mutationResolver) CreateAccountBuyer(ctx context.Context, in AccountBuyerInput) (*AccountBuyer, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	userAuth, err := m.GetUserContext(ctx)
	if err != nil {
		return nil, err
	}

	data := &pb.PostAccountBuyerRequest{}
	req := MapGraphQLInputToRequest(in, TypeCreate)
	switch v := req.(type) {
	case *pb.PostAccountBuyerRequest:
		data = v
	default:
		return nil, ErrInvalidInfo
	}

	data.Id = userAuth.ID
	a, err := m.server.accountClient.PostAccountBuyer(ctx, data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &AccountBuyer{
		ID:        a.ID,
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Phone:     a.Phone,
		Address:   a.Address,
		Orders:    []*Order{},
	}, nil
}

func (m *mutationResolver) CreateAccountSeller(ctx context.Context, in AccountSellerInput) (*AccountSeller, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	userAuth, err := m.GetUserContext(ctx)
	if err != nil {
		return nil, err
	}

	data := &pb.PostAccountSellerRequest{}
	req := MapGraphQLInputToRequest(in, TypeCreate)
	switch v := req.(type) {
	case *pb.PostAccountSellerRequest:
		data = v
	default:
		return nil, ErrInvalidInfo
	}

	data.Id = userAuth.ID
	a, err := m.server.accountClient.PostAccountSeller(ctx, data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &AccountSeller{
		ID:        a.ID,
		StoreName: a.StoreName,
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Phone:     a.Phone,
		Address:   a.Address,
		Products:  []*Product{},
	}, nil
}

func (m *mutationResolver) UpdateAccountBuyer(ctx context.Context, in AccountBuyerInput, id string) (*AccountBuyer, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	data := &pb.AccountBuyer{}
	req := MapGraphQLInputToRequest(in, TypeUpdate)
	switch v := req.(type) {
	case *pb.AccountBuyer:
		v.Id = id
		data = v
	default:
		return nil, ErrInvalidInfo
	}

	a, err := m.server.accountClient.UpdateAccountBuyer(ctx, data)
	if err != nil {
		return nil, err
	}

	return &AccountBuyer{
		ID:        a.Id,
		FirstName: a.BaseInfo.FirstName,
		LastName:  a.BaseInfo.LastName,
		Phone:     a.BaseInfo.Phone,
		Address:   a.BaseInfo.Address,
		Orders:    []*Order{},
	}, nil
}

func (m *mutationResolver) UpdateAccountSeller(ctx context.Context, in AccountSellerInput, id string) (*AccountSeller, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	data := &pb.AccountSeller{}
	req := MapGraphQLInputToRequest(in, TypeUpdate)
	switch v := req.(type) {
	case *pb.AccountSeller:
		v.Id = id
		data = v
	default:
		return nil, ErrInvalidInfo
	}

	log.Println("data: ", data)
	a, err := m.server.accountClient.UpdateAccountSeller(ctx, data)
	if err != nil {
		return nil, err
	}

	return &AccountSeller{
		ID:        a.Id,
		FirstName: a.BaseInfo.FirstName,
		LastName:  a.BaseInfo.LastName,
		Phone:     a.BaseInfo.Phone,
		Address:   a.BaseInfo.Address,
		Products:  []*Product{},
	}, nil
}

func (m *mutationResolver) LoginUser(ctx context.Context, email, password string) (LoginResult, error) {
	w, ok := ctx.Value(responseWriterKey).(http.ResponseWriter)
	if !ok {
		return nil, &gqlerror.Error{Message: "response writer not found in context"}
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	user, err := m.server.authClient.LoginUser(ctx, &auth_pb.LoginUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	cookieToken := &http.Cookie{
		Name:     "token",
		Value:    user.TokenResponse.GetToken(),
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	cookieRefreshToken := &http.Cookie{
		Name:     "refresh_token",
		Value:    user.TokenResponse.GetRefreshToken(),
		Expires:  time.Now().Add(240 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	http.SetCookie(w, cookieToken)
	http.SetCookie(w, cookieRefreshToken)

	// access role from context
	role := MapIntToRole(user.Role)
	if role == RoleTypeSeller {
		return &AccountSeller{
			Email: email,
		}, nil
	}

	// Default: Return an AccountBuyer
	return &AccountBuyer{
		Email: email,
	}, nil
}

func (m *mutationResolver) RefrehToken(ctx context.Context, token string) (*RefreshToken, error) {
	w, ok := ctx.Value(responseWriterKey).(http.ResponseWriter)
	if !ok {
		return nil, &gqlerror.Error{Message: "response writer not found in context"}
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	in := &auth_pb.RefreshTokenRequest{
		RefreshToken: token,
	}
	tokenPair, err := m.server.authClient.RefreshTokenUser(ctx, in)
	if err != nil {
		return nil, err
	}

	cookieToken := &http.Cookie{
		Name:     "token",
		Value:    tokenPair.GetToken(),
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	cookieRefreshToken := &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.GetRefreshToken(),
		Expires:  time.Now().Add(240 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	http.SetCookie(w, cookieToken)
	http.SetCookie(w, cookieRefreshToken)

	return &RefreshToken{}, nil
}

func (m *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	p, err := m.server.catalogClient.PostProduct(ctx, in.Name, in.Description, in.SellerID, in.Price, uint32(in.Quantity))
	if err != nil {
		return nil, err
	}

	log.Println(p.SellerID)
	return &Product{
		ID:          p.ID,
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Quantity:    in.Quantity,
		SellerID:    p.SellerID,
	}, nil
}

func (m *mutationResolver) UpdateProduct(ctx context.Context, in ProductInput, id string) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// edit function update product
	_, err := m.server.catalogClient.UpdateProduct(ctx, id, in.Name, in.Description, in.Price, uint32(in.Quantity))
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          id,
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Quantity:    in.Quantity,
	}, nil
}

func (m *mutationResolver) DeleteProduct(ctx context.Context, id string) (string, error) {
	// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	// defer cancel()
	return id, nil
}

func (m *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var products []order.OrderedProduct
	for _, p := range in.Products {
		if p.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}
		products = append(products, order.OrderedProduct{
			ID:       p.ProductID,
			Quantity: uint32(p.Quantity),
		})
	}

	order, err := m.server.orderClient.PostOrder(ctx, in.AccountID, products)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Order{
		ID:         order.ID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
	}, nil
}

func (m *mutationResolver) DeleteOrder(ctx context.Context, id string) (string, error) {
	// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	// defer cancel()
	return id, nil
}

func (m *mutationResolver) GetUserContext(ctx context.Context) (UserAuth, error) {
	u, ok := ctx.Value(userCtxKey).(UserAuth)
	if !ok {
		return u, &gqlerror.Error{Message: "user not found in context"}
	}

	return u, nil
}
