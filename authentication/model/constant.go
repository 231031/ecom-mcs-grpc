package model

type RoleType int

const (
	BUYER RoleType = iota
	SELLER
)

var RoleTypeLabel = map[RoleType]string{
	BUYER:  "buyer",
	SELLER: "seller",
}

func (c RoleType) String() string {
	return RoleTypeLabel[c]
}

func (c RoleType) IsValid() bool {
	_, ok := RoleTypeLabel[c]
	return ok
}
