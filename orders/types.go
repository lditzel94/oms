package main

import (
	"context"
	pb "github.com/lditzel94/oms/commons/api"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest, []*pb.Item) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
}

type OrdersStore interface {
	Create(context.Context, *pb.CreateOrderRequest, []*pb.Item) (string, error)
	Get(ctx context.Context, id, customerID string) (*pb.Order, error)
}
