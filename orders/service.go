package main

import (
	"context"
	"github.com/lditzel94/oms/commons"
	pb "github.com/lditzel94/oms/commons/api"
	"log"
)

type service struct {
	store OrdersStore
}

func NewService(store *store) *service {
	return &service{store: store}
}

func (s *service) CreateOrder(ctx context.Context) error {
	return nil
}

func (s *service) ValidateOrder(ctx context.Context, request *pb.CreateOrderRequest) error {
	if len(request.Items) == 0 {
		return commons.ErrorNoItems
	}

	mergedItems := mergedItemsQuantities(request.Items)
	log.Print(mergedItems)

	//validate with the stock service
	return nil
}

func mergedItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)
	for _, item := range items {
		found := false

		for _, finalItem := range merged {
			if finalItem.ID == item.ID {
				finalItem.Quantity += item.Quantity
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, item)
		}
	}

	return merged
}
