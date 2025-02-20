package main

import (
	"context"
	"github.com/lditzel94/oms/commons/api"
	"github.com/lditzel94/oms/payments/processor/inmem"
	"testing"
)

func TestService(t *testing.T) {
	processor := inmem.NewInmem()
	svc := NewService(processor)

	t.Run("should create a payment link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &api.Order{})
		if err != nil {
			t.Errorf("CreatePayment() error = %v, want nil", err)
		}

		if link == "" {
			t.Error("CreatePayment() link is empty")
		}
	})
}
