package processor

import pb "github.com/lditzel94/oms/commons/api"

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
