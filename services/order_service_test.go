package services

import (
	"context"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"orderservice/models"
	pb "orderservice/proto"
	"regexp"
	"testing"
)

func setupOrderServiceTest() (sqlmock.Sqlmock, *OrderService) {
	mockdb, mock, _ := sqlmock.New()
	dialector := postgres.New(
		postgres.Config{
			DriverName: "postgres",
			Conn:       mockdb,
		})
	gormDB, _ := gorm.Open(dialector, &gorm.Config{})
	service := OrderService{Database: gormDB}

	return mock, &service
}

func TestCreateOrder(t *testing.T) {
	mock, service := setupOrderServiceTest()

	req := &pb.CreateOrderRequest{
		User: &pb.User{
			Id: 123,
			Location: &pb.Location{
				Latitude:  45.0,
				Longitude: -75.0,
			},
		},
		Restaurant: &pb.Restaurant{
			Id: 456,
			Location: &pb.Location{
				Latitude:  46.0,
				Longitude: -76.0,
			},
		},
		Items: []*pb.MenuItem{
			{
				Id:    789,
				Name:  "Pizza",
				Price: 10.99,
			},
		},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

	ctx := context.Background()
	resp, err := service.Create(ctx, req)
	var expectedid int64 = 0
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedid, resp.OrderId)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestCreateDeliveryForOrder(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := "http://localhost:8090/fullfilment-management/deliveries"

	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

			var payload map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&payload)
			assert.NoError(t, err)

			var expectedid float64 = 1
			assert.Equal(t, expectedid, payload["orderId"])
			assert.Equal(t, 100.0, payload["totalPrice"].(float64))

			return httpmock.NewStringResponse(201, ""), nil
		},
	)

	order := models.Order{
		ID:         1,
		TotalPrice: 100,
		Restaurant: models.Restaurant{
			ID:                 1,
			RestaurantUniqueID: 12,
			Location:           models.Location{Latitude: 45.0, Longitude: -75.0},
		},
		User: models.User{Location: models.Location{Latitude: 46.0, Longitude: -76.0}},
	}

	createDeliveryForOrder(order)

	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info["POST "+url])
}
