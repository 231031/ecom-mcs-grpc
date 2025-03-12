package catalog

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func mapProductResponse(productResp []productResp, products *[]Product) {
	for _, p := range productResp {
		*products = append(*products, Product{
			ID:          p.ID,
			Name:        p.Source.Name,
			Description: p.Source.Description,
			Price:       p.Source.Price,
		})
	}
}

func marshalJsonID(ids []string) ([]byte, error) {
	idDocs := []map[string]string{}
	for _, id := range ids {
		idJson := map[string]string{
			"_id": id,
		}
		idDocs = append(idDocs, idJson)
	}

	docs := map[string]interface{}{
		"docs": idDocs,
	}

	docsJson, err := json.Marshal(docs)
	if err != nil {
		return nil, err
	}

	return docsJson, nil
}

func createQueryCheck(p Product) string {
	priceProduct := strconv.FormatFloat(p.Price, 'f', -1, 64)
	return fmt.Sprintf("name:%s AND price:%s AND description:%s", p.Name, priceProduct, p.Description)
}

func convertProductToMap(p Product) (map[string]interface{}, error) {
	// Marshall the Product struct to JSON
	productJSON, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON into a map
	var result map[string]interface{}
	err = json.Unmarshal(productJSON, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
