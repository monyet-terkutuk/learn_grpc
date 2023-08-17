package services

import (
	"context"
	"learn_grpc/pb/pagination"
	productPB "learn_grpc/pb/product"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ProductService struct {
	productPB.UnimplementedProductServiceServer
	DB *gorm.DB
}

func (p *ProductService) GetProducts(context.Context, *productPB.Empty) (*productPB.Products, error) {
	var products []*productPB.Product

	rows, err := p.DB.Table("products AS p").Joins("LEFT JOIN categories AS c ON c.id = p.category_id").Select("p.id", "p.name", "p.price", "p.stock", "c.id AS category_id", "c.name AS category_name").Rows()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var product productPB.Product
		var category productPB.Category

		if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Stock, &category.Id, &category.Name); err != nil {
			log.Fatalf("Failed to get row data %v", err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}

		product.Category = &category
		products = append(products, &product)
	}

	response := &productPB.Products{
		Pagination: &pagination.Pagination{
			Total:       5,
			PerPage:     3,
			CurrentPage: 2,
			LastPage:    2,
		},
		Data: products,
	}

	return response, nil
}
