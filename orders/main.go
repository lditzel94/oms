package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lditzel94/oms/commons"
	"github.com/lditzel94/oms/commons/broker"
	"github.com/lditzel94/oms/commons/discovery"
	"github.com/lditzel94/oms/commons/discovery/consul"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var (
	serviceName = "orders"
	grpcAddr    = commons.EnvString("GRPC_ADDR", "localhost:2000")
	consulAddr  = commons.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser    = commons.EnvString("RABBITMQ_USER", "guest")
	amqpPass    = commons.EnvString("RABBITMQ_PASS", "guest")
	amqpHost    = commons.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort    = commons.EnvString("RABBITMQ_PORT", "5672")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("health check failed: ", err)
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, serviceName, instanceID)

	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer l.Close()

	store := NewStore()
	service := NewService(store)
	NewGRPCHandler(grpcServer, service, ch)

	service.CreateOrder(context.Background())

	log.Println("GRPC server listening on", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
