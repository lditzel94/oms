package main

import (
	"context"
	pb "github.com/lditzel94/oms/commons/api"
)

type OrdersService interface {
	CreateOrder(ctx context.Context) error
	ValidateOrder(ctx context.Context, request *pb.CreateOrderRequest) error
}

type OrdersStore interface {
	Create(ctx context.Context) error
}
