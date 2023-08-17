package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/go_grpc"))
	if err != nil {
		log.Fatalf("Database connection failed %v", err.Error())
	}
	fmt.Println("Database connected")
	return db
}
