package main

import (
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net"
	"orderservice/models"
	pb "orderservice/proto"
	"orderservice/services"
)

const (
	port = ":9001"
)

func databaseConn() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=orderservice port=5433 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("error connecting to database: %s", err)
	}

	err = db.AutoMigrate(&models.Order{}, &models.MenuItem{}, &models.User{})

	if err != nil {
		return nil
	}

	db.Logger.LogMode(logger.Info)

	return db
}

func main() {
	db := databaseConn()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error listening to port: %s", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterOrderServiceServer(grpcServer, &services.OrderService{Database: db})

	log.Println("Listening to port", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to Server: %s", err)
	}
}
