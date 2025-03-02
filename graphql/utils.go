package graphql

import (
	"github.com/231031/ecom-mcs-grpc/account/pb"
)

type SelectionType string

const (
	TypeCreate SelectionType = "CREATE"
	TypeUpdate SelectionType = "UPDATE"
)

func MapGraphQLInputToRequest(data any, selType SelectionType) any {

	switch in := data.(type) {
	case AccountBuyerInput:
		if selType == "CREATE" {
			return &pb.PostAccountBuyerRequest{
				BaseInfo: &pb.BaseInfo{
					Email:     in.BaseInfo.Email,
					FirstName: in.BaseInfo.FirstName,
					LastName:  in.BaseInfo.LastName,
					Phone:     in.BaseInfo.Phone,
					Address:   in.BaseInfo.Address,
				},
				Password: in.Password,
			}
		} else if selType == "UPDATE" {
			return &pb.AccountBuyer{
				BaseInfo: &pb.BaseInfo{
					Email:     in.BaseInfo.Email,
					FirstName: in.BaseInfo.FirstName,
					LastName:  in.BaseInfo.LastName,
					Phone:     in.BaseInfo.Phone,
					Address:   in.BaseInfo.Address,
				},
			}
		}
	case AccountSellerInput:
		if selType == "CREATE" {
			return &pb.PostAccountSellerRequest{
				BaseInfo: &pb.BaseInfo{
					Email:     in.BaseInfo.Email,
					FirstName: in.BaseInfo.FirstName,
					LastName:  in.BaseInfo.LastName,
					Phone:     in.BaseInfo.Phone,
					Address:   in.BaseInfo.Address,
				},
				StoreName: in.StoreName,
				Password:  in.Password,
			}
		} else if selType == "UPDATE" {
			return &pb.AccountSeller{
				BaseInfo: &pb.BaseInfo{
					Email:     in.BaseInfo.Email,
					FirstName: in.BaseInfo.FirstName,
					LastName:  in.BaseInfo.LastName,
					Phone:     in.BaseInfo.Phone,
					Address:   in.BaseInfo.Address,
				},
				StoreName: in.StoreName,
			}
		}
	}
	return nil
}
