package main

import (
	"errors"
	"github.com/lditzel94/oms/commons"
	pb "github.com/lditzel94/oms/commons/api"
	"github.com/lditzel94/oms/gateway/gateway"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type handler struct {
	gateway gateway.OrdersGateway
}

func NewHandler(gateway gateway.OrdersGateway) *handler {
	return &handler{gateway}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	var items []*pb.ItemsWithQuantity
	if err := commons.ReadJSON(r, &items); err != nil {
		commons.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateItems(items); err != nil {
		commons.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.gateway.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})

	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() == codes.InvalidArgument {
			commons.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}

		commons.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	commons.WriteJSON(w, http.StatusCreated, order)
}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return commons.ErrorNoItems
	}

	for _, item := range items {
		if item.ID == "" {
			return errors.New("item must have an ID")
		}

		if item.Quantity <= 0 {
			return errors.New("item must have a valid Quantity")
		}
	}

	return nil
}
