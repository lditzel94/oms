package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	common "github.com/lditzel94/oms/commons"
	"github.com/lditzel94/oms/commons/broker"
	"github.com/lditzel94/oms/commons/discovery"
	"github.com/lditzel94/oms/commons/discovery/consul"
	stripeProcessor "github.com/lditzel94/oms/payments/processor/stripe"
	"github.com/stripe/stripe-go/v78"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var (
	serviceName = "payment"
	amqpUser    = common.EnvString("RABBITMQ_USER", "guest")
	amqpPass    = common.EnvString("RABBITMQ_PASS", "guest")
	amqpHost    = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort    = common.EnvString("RABBITMQ_PORT", "5672")
	grpcAddr    = common.EnvString("GRPC_ADDRESS", "localhost:2001")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	stripeKey   = common.EnvString("STRIPE_KEY", "")
	//httpAddr             = common.EnvString("HTTP_ADDR", "localhost:8081")
	//endpointStripeSecret = common.EnvString("STRIPE_ENDPOINT_SECRET", "whsec_...")
	//jaegerAddr           = common.EnvString("JAEGER_ADDR", "localhost:4318")
)

func main() {
	//if err := common.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr); err != nil {
	//	log.Fatal("failed to set global tracer")
	//}

	// Register consul
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	instanceID := discovery.GenerateInstanceID(serviceName)

	ctx := context.Background()
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatalf("failed to health check %v", err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	// stripe setup
	stripe.Key = stripeKey

	// Broker connection
	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	stripeProcessor := stripeProcessor.NewProcessor()
	svc := NewService(stripeProcessor)

	amqpConsumer := NewConsumer(svc)
	go amqpConsumer.Listen(ch)

	//svcWithTelemetry := NewTelemetryMiddleware(svc)
	//gateway := gateway.NewGateway(registry)

	// http server
	//mux := http.NewServeMux()

	//httpServer := NewPaymentHTTPHandler(ch)
	//httpServer.registerRoutes(mux)

	//go func() {
	//	log.Printf("Starting HTTP server at %s", httpAddr)
	//	if err := http.ListenAndServe(httpAddr, mux); err != nil {
	//		log.Fatal("failed to start http server")
	//	}
	//}()

	// gRPC server
	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	log.Println("GRPC Server Started at ", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
