package services

import (
	"context"
	productPB "learn_grpc/pb/product"

	"gorm.io/gorm"
)

type ProductService struct {
	productPB.UnimplementedProductServiceServer
	DB *gorm.DB
}

func (p *ProductService) GetProducts(context.Context, *productPB.Empty) (*productPB.Products, error) {
	var products []*productPB.Products

	rows, err := p.DB.Table()

	return products, nil
}
