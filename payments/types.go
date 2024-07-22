package main

import (
	"context"
	pb "github.com/lditzel94/oms/commons/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
