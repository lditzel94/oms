package main

import (
	"context"
	pb "github.com/lditzel94/oms/commons/api"
)

type OrdersService interface {
	CreateOrder(ctx context.Context, request *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrder(ctx context.Context, request *pb.CreateOrderRequest) ([]*pb.Item, error)
}

type OrdersStore interface {
	Create(ctx context.Context) error
}
