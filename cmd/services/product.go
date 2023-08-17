package services

import (
	"context"
	"learn_grpc/cmd/helpers"
	pagination "learn_grpc/pb/pagination"
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

func (p *ProductService) GetProducts(ctx context.Context, pageParam *productPB.Page) (*productPB.Products, error) {
	var page int64 = 1
	if pageParam.GetPage() != 0 {
		page = pageParam.GetPage()
	}

	var pagination pagination.Pagination
	var products []*productPB.Product

	sql := p.DB.Table("products AS p").Joins("LEFT JOIN categories AS c ON c.id = p.category_id").Select("p.id", "p.name", "p.price", "p.stock", "c.id AS category_id", "c.name AS category_name")

	offset, limit := helpers.Pagination(sql, page, &pagination)
	rows, err := sql.Offset(int(offset)).Limit(int(limit)).Rows()

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
		Pagination: &pagination,
		Data:       products,
	}

	return response, nil
}

func (p *ProductService) GetProduct(ctx context.Context, id *productPB.Id) (*productPB.Product, error) {
	row := p.DB.Table("products AS p").Joins("LEFT JOIN categories AS c ON c.id = p.category_id").Select("p.id", "p.name", "p.price", "p.stock", "c.id AS category_id", "c.name AS category_name").Where("p.id = ?", id.GetId()).Row()

	var product productPB.Product
	var category productPB.Category

	if err := row.Scan(&product.Id, &product.Name, &product.Price, &product.Stock, &category.Id, &category.Name); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	product.Category = &category

	return &product, nil
}

func (p *ProductService) CreateProduct(ctx context.Context, productData *productPB.Product) (*productPB.Id, error) {
	var Response productPB.Id

	err := p.DB.Transaction(func(tx *gorm.DB) error {
		category := productPB.Category{
			Id:   0,
			Name: productData.GetCategory().Name,
		}

		if err := tx.Table("categories").Where("LCASE(name) = ?", category.Name).FirstOrCreate(&category).Error; err != nil {
			return err
		}

		product := struct {
			ID         uint64
			Name       string
			Price      float64
			Stock      uint32
			CategoryID uint32
		}{
			ID:         productData.GetId(),
			Name:       productData.GetName(),
			Price:      productData.GetPrice(),
			Stock:      productData.GetStock(),
			CategoryID: category.Id,
		}

		if err := tx.Table("products").Create(&product).Error; err != nil {
			return err
		}

		Response.Id = product.ID
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &Response, nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, productData *productPB.Product) (*productPB.Status, error) {
	var Response productPB.Status

	err := p.DB.Transaction(func(tx *gorm.DB) error {
		category := productPB.Category{
			Id:   0,
			Name: productData.GetCategory().Name,
		}

		if err := tx.Table("categories").Where("LCASE(name) = ?", category.Name).FirstOrCreate(&category).Error; err != nil {
			return err
		}

		product := struct {
			ID         uint64
			Name       string
			Price      float64
			Stock      uint32
			CategoryID uint32
		}{
			ID:         productData.GetId(),
			Name:       productData.GetName(),
			Price:      productData.GetPrice(),
			Stock:      productData.GetStock(),
			CategoryID: category.Id,
		}

		if err := tx.Table("products").Where("id = ?", product.ID).Updates(&product).Error; err != nil {
			return err
		}

		Response.Status = 1
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &Response, nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, ID *productPB.Id) (*productPB.Status, error) {
	var response productPB.Status

	if err := p.DB.Table("products").Where("id = ?", ID.GetId()).Delete(nil).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response.Status = 1

	return &response, nil
}
