package main

import (
	"context"
	pb "github.com/lditzel94/oms/commons/api"
	"github.com/lditzel94/oms/payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
	//gateway gateway.OrdersGateway
}

func NewService(processor processor.PaymentProcessor) *service {
	return &service{processor: processor}
}

func (s *service) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	link, err := s.processor.CreatePaymentLink(o)
	if err != nil {
		return "", err
	}

	// update order with the link
	//err = s.gateway.UpdateOrderAfterPaymentLink(ctx, o.ID, link)
	//if err != nil {
	//	return "", err
	//}

	return link, nil
}
