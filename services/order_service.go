package services

import (
	"context"
	"gorm.io/gorm"
	"orderservice/models"
	pb "orderservice/proto"
)

type OrderService struct {
	Database *gorm.DB
	pb.UnimplementedOrderServiceServer
}

func (server *OrderService) Create(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	var items []models.MenuItem
	var totalPrice float32 = 0

	for _, item := range req.Items {
		items = append(items, models.MenuItem{
			ItemUniqueID: item.GetId(),
			Name:         item.GetName(),
			Price:        item.GetPrice(),
		})
		totalPrice += item.GetPrice()
	}

	user := models.User{
		UserUniqueID: req.User.GetId(),
		Location: models.Location{
			Latitude:  req.User.Location.GetLatitude(),
			Longitude: req.User.Location.GetLongitude(),
		},
	}
	restaurant := models.Restaurant{
		RestaurantUniqueID: req.Restaurant.GetId(),
		Location: models.Location{
			Latitude:  req.Restaurant.Location.GetLatitude(),
			Longitude: req.Restaurant.Location.GetLongitude(),
		},
	}

	order := models.Order{
		Restaurant: restaurant,
		User:       user,
		Status:     "BOOKED",
		Items:      items,
		TotalPrice: totalPrice,
	}

	server.Database.Create(&order)

	return &pb.CreateOrderResponse{
		OrderId: order.ID,
		Restaurant: &pb.Restaurant{
			Id: order.Restaurant.RestaurantUniqueID,
			Location: &pb.Location{
				Latitude:  order.Restaurant.Location.Latitude,
				Longitude: order.Restaurant.Location.Longitude,
			},
		},
		Items:  req.Items,
		Status: order.Status,
		User: &pb.User{
			Id: order.User.UserUniqueID,
			Location: &pb.Location{
				Latitude:  order.User.Latitude,
				Longitude: order.User.Longitude,
			},
		},
		TotalPrice: order.TotalPrice,
	}, nil
}
