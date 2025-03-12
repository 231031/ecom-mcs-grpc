package catalog

type mGetResp struct {
	Hits []productResp `json:"docs"`
}

type listsProductResp struct {
	Hits hitsArray `json:"hits"`
}

type hitsArray struct {
	Hits []productResp `json:"hits"`
}

type productResp struct {
	ID     string          `json:"_id"`
	Found  bool            `json:"found"`
	Source productDocument `json:"_source"`
}

type productDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    uint32  `json:"quantity"`
	SellerID    string  `json:"seller_id"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    uint32  `json:"quantity"`
	SellerID    string  `json:"seller_id"`
}
