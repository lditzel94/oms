package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lditzel94/oms/commons"
	"github.com/lditzel94/oms/commons/discovery"
	"github.com/lditzel94/oms/commons/discovery/consul"
	"github.com/lditzel94/oms/gateway/gateway"
	"log"
	"net/http"
	"time"
)

var (
	serviceName = "gateway"
	httpAddr    = commons.EnvString("HTTP_ADDR", ":3000")
	consulAddr  = commons.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
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

	mux := http.NewServeMux()

	ordersGw := gateway.NewGRPCGateway(registry)
	handler := NewHandler(ordersGw)
	handler.registerRoutes(mux)

	log.Printf("Starting server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start http server")
	}
}
