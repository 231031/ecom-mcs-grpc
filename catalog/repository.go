package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

var (
	ErrNotFound     = errors.New("product not found")
	ErrAlreadyExist = errors.New("product is already registered")
	ErrPutProduct   = errors.New("falied to put product")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
	UpdateQuantity(ctx context.Context, ids []string, quantity []uint32) error
	UpdateProduct(ctx context.Context, p map[string]interface{}) error
}

type elasticRepository struct {
	client *elasticsearch.Client
}

func NewElasticRepository(url, username, password string) (Repository, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			url,
		},
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return &elasticRepository{client: client}, nil
}

func (r *elasticRepository) Close() {
}

func (r *elasticRepository) PutProduct(ctx context.Context, p Product) error {
	query := createQueryCheck(p)
	pFound, err := r.SearchProducts(ctx, query, 0, 1)
	if err != nil {
		return err
	}

	if len(pFound) != 0 {
		return ErrAlreadyExist
	}

	product := productDocument{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Quantity:    p.Quantity,
		SellerID:    p.SellerID,
	}
	productJson, err := json.Marshal(product)
	if err != nil {
		return err
	}

	resp, err := r.client.Index(
		"products",
		bytes.NewReader(productJson),
		r.client.Index.WithDocumentID(p.ID),
		r.client.Index.WithContext(ctx),
	)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		log.Println(resp.String())
		return ErrPutProduct
	}

	return nil
}

func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	resp, err := r.client.Get(
		"products",
		id,
		r.client.Get.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	p := productResp{}
	err = json.Unmarshal(body, &p)
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          id,
		Name:        p.Source.Name,
		Description: p.Source.Description,
		Price:       p.Source.Price,
		Quantity:    p.Source.Quantity,
	}, nil
}

func (r *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	resp, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("products"),
		r.client.Search.WithSize(int(take)),
		r.client.Search.WithFrom(int(skip)),
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	listResp := listsProductResp{}
	err = json.Unmarshal(body, &listResp)
	if err != nil {
		return nil, err
	}

	products := []Product{}
	mapProductResponse(listResp.Hits.Hits, &products)
	return products, nil
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	idsJson, err := marshalJsonID(ids)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Mget(
		bytes.NewReader(idsJson),
		r.client.Mget.WithContext(ctx),
		r.client.Mget.WithIndex("products"),
	)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	listResp := mGetResp{}
	err = json.Unmarshal(body, &listResp)
	if err != nil {
		return nil, err
	}

	products := []Product{}
	for _, p := range listResp.Hits {
		if p.Found {
			products = append(products, Product{
				ID:          p.ID,
				Name:        p.Source.Name,
				Description: p.Source.Description,
				Price:       p.Source.Price,
				Quantity:    p.Source.Quantity,
			})
		}
	}

	return products, nil
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	resp, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("products"),
		r.client.Search.WithQuery(query),
		r.client.Search.WithFrom(int(skip)),
		r.client.Search.WithSize(int(take)),
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	listResp := listsProductResp{}
	err = json.Unmarshal(body, &listResp)
	if err != nil {
		return nil, err
	}

	products := []Product{}
	mapProductResponse(listResp.Hits.Hits, &products)
	return products, nil
}

func (r *elasticRepository) UpdateQuantity(ctx context.Context, ids []string, quantity []uint32) error {
	var builder strings.Builder
	for i := range ids {
		builder.WriteString(fmt.Sprintf(`{ "update": { "_index": "products", "_id": "%s" } }%s`, ids[i], "\n"))
		builder.WriteString(fmt.Sprintf(`{ "doc" : {"quantity" : "%s"} }%s`, strconv.Itoa(int(quantity[i])), "\n"))
	}
	body := builder.String()
	body += "\n"

	if body != "" {
		resp, err := r.client.Bulk(bytes.NewReader([]byte(body)))
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			log.Println(resp.String())
		}
		return nil
	}

	return nil
}

func (r *elasticRepository) UpdateProduct(ctx context.Context, p map[string]interface{}) error {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(`{ "update": { "_index": "products", "_id": "%s" } }%s`, p["ID"], "\n"))

	builder.WriteString(`{ "doc" : {`)
	first := true
	for k, v := range p {
		if k != "ID" && v != "" {
			if !first {
				builder.WriteString(", ")
			}
			builder.WriteString(fmt.Sprintf(`"%s" : "%s"`, k, v))
			first = false
		}
	}
	builder.WriteString(`} }` + "\n")

	body := builder.String()
	body += "\n"

	if body != "" {
		resp, err := r.client.Bulk(bytes.NewReader([]byte(body)))
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			log.Println(resp.String())
		}
		return nil
	}

	return nil
}
