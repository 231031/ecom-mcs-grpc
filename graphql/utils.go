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
					FirstName: in.BaseInfo.FirstName,
					LastName:  in.BaseInfo.LastName,
					Phone:     in.BaseInfo.Phone,
					Address:   in.BaseInfo.Address,
				},
			}
		} else if selType == "UPDATE" {
			return &pb.AccountBuyer{
				BaseInfo: &pb.BaseInfo{
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
					FirstName: in.BaseInfo.FirstName,
					LastName:  in.BaseInfo.LastName,
					Phone:     in.BaseInfo.Phone,
					Address:   in.BaseInfo.Address,
				},
				StoreName: in.StoreName,
			}
		} else if selType == "UPDATE" {
			return &pb.AccountSeller{
				BaseInfo: &pb.BaseInfo{
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

func MapRoleToInt(role RoleType) int32 {
	mapRole := map[RoleType]int32{
		RoleTypeBuyer:  0,
		RoleTypeSeller: 1,
	}

	return mapRole[role]
}

func MapIntToRole(roleNum int32) RoleType {
	mapRole := map[int32]RoleType{
		0: RoleTypeBuyer,
		1: RoleTypeSeller,
	}
	return mapRole[roleNum]
}
