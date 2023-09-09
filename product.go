package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/google/uuid"
)

var ProductIndex = "product"

type Product struct {
	// Key
	ID       uuid.UUID `json:"id"`
	SellerID uuid.UUID `json:"seller_id"`

	// Attr
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`

	// Metadata
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ESProduct struct {
	ID     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source Product `json:"_source"`
}

func (p *Product) Marshal() ([]byte, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal product: %s", err)
	}

	return b, nil
}

func (p *Product) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, &p)
	if err != nil {
		return fmt.Errorf("failed to unmarshal product: %s", err)
	}

	return nil
}

func (p *Product) Index(ctx context.Context, es *elasticsearch.TypedClient) error {
	_, err := es.Index(ProductIndex).
		Id(p.ID.String()).
		Request(p).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to index product with id %s: %s", p.ID, err)
	}

	return nil
}

func (p *Product) GetByID(ctx context.Context, id uuid.UUID, es *elasticsearch.TypedClient) error {
	resp, err := es.Get(ProductIndex, id.String()).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to get document with id %s: %s", id, err)
	}

	b, err := json.Marshal(resp.Source_)
	if err != nil {
		return fmt.Errorf("failed to marshal document result: %s", err)
	}

	if err = p.Unmarshal(b); err != nil {
		return fmt.Errorf("failed to unmarshal document: %s", err)
	}

	return nil
}

func (p *Product) Search(ctx context.Context, es *elasticsearch.TypedClient) ([]ESProduct, error) {
	resp, err := es.Search().
		Index(ProductIndex).
		Request(&search.Request{
			Query: &types.Query{Bool: &types.BoolQuery{
				Must: []types.Query{
					{Term: map[string]types.TermQuery{
						"name": {
							Value: "test",
						},
					}},
				},
			}},
		}).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to search document: %s", err)
	}

	if resp.Hits.Total.Value == 0 {
		return nil, fmt.Errorf("document not found")
	}

	var products []ESProduct
	for _, v := range resp.Hits.Hits {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal document instance: %v", err)
		}

		var p ESProduct
		err = json.Unmarshal(b, &p)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal document instance: %v", err)
		}

		products = append(products, p)
	}

	return products, nil
}
