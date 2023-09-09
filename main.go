package main

import (
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	conf, err := NewConfig(".env")
	if err != nil {
		panic(err)
	}

	es, err := NewElasticSearch(conf.CloudID, conf.APIKey)
	if err != nil {
		panic(err)
	}

	// p := Product{
	// 	ID:          uuid.New(),
	// 	SellerID:    uuid.New(),
	// 	Name:        "test from go",
	// 	Description: "this is a test from go lang",
	// 	Price:       10000,
	// 	CreatedAt:   time.Now(),
	// }

	// err = p.Index(ctx, es)
	// if err != nil {
	// 	panic(err)
	// }

	p := Product{}
	products, err := p.Search(ctx, es)
	if err != nil {
		panic(err)
	}

	log.Print(products[0].Source.ID, " && ", products[0].ID)

	err = p.GetByID(ctx, products[0].Source.ID, es)
	if err != nil {
		panic(err)
	}

	log.Print(p.ID)
}
