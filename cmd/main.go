package main

import (
	"context"
	"log"
	"net"

	"github.com/justIGreK/MoneyKeeper-Transaction/cmd/handler"
	"github.com/justIGreK/MoneyKeeper-Transaction/internal/repository"
	"github.com/justIGreK/MoneyKeeper-Transaction/internal/service"
	"github.com/justIGreK/MoneyKeeper-Transaction/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	user, err := client.NewUserClient("localhost:50052")
	if err != nil {
		log.Fatal(err)
	}
	db := repository.CreateMongoClient(ctx)
	txRepo := repository.NewTransactionRepository(db)
	txSRV := service.NewTransactionService(txRepo, user)
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	handler := handler.NewHandler(grpcServer, txSRV)
	handler.RegisterServices()
	reflection.Register(grpcServer)

	log.Printf("Starting gRPC server on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
