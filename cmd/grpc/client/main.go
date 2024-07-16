package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "gitlab.ozon.dev/r_gabdullin/homework-1/pb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	acceptOrder(ctx, client)
	getOrders(ctx, client)
	deliverOrder(ctx, client)
	getReturns(ctx, client)
	returnOrder(ctx, client)
	acceptReturn(ctx, client)
}

func acceptOrder(ctx context.Context, client pb.OrderServiceClient) {
	req := &pb.AcceptOrderRequest{
		User:      1,
		Order:     101,
		Weight:    2,
		BasePrice: 1000,
		Expire:    "2024-12-20T12",
		Wrapper:   "pack",
	}
	_, err := client.AcceptOrder(ctx, req)
	if err != nil {
		log.Fatalf("could not accept order: %v", err)
	}
	fmt.Println("Order accepted")
}

func acceptReturn(ctx context.Context, client pb.OrderServiceClient) {
	req := &pb.AcceptReturnRequest{
		User:  1,
		Order: 101,
	}
	_, err := client.AcceptReturn(ctx, req)
	if err != nil {
		log.Fatalf("could not accept return: %v", err)
	}
	fmt.Println("Return accepted")
}

func deliverOrder(ctx context.Context, client pb.OrderServiceClient) {
	req := &pb.DeliverOrderRequest{
		Orders: []int32{101},
	}
	_, err := client.DeliverOrder(ctx, req)
	if err != nil {
		log.Fatalf("could not deliver orders: %v", err)
	}
	fmt.Println("Orders delivered")
}

func getOrders(ctx context.Context, client pb.OrderServiceClient) {
	req := &pb.GetOrdersRequest{
		User:  1,
		Count: 10,
	}
	resp, err := client.GetOrders(ctx, req)
	if err != nil {
		log.Fatalf("could not get orders: %v", err)
	}
	fmt.Println("Orders:", resp.GetOrders())
}

func getReturns(ctx context.Context, client pb.OrderServiceClient) {
	req := &pb.GetReturnsRequest{
		Offset: 0,
		Limit:  10,
	}
	resp, err := client.GetReturns(ctx, req)
	if err != nil {
		log.Fatalf("could not get returns: %v", err)
	}
	fmt.Println("Returns:", resp.GetReturns())
}

func returnOrder(ctx context.Context, client pb.OrderServiceClient) {
	req := &pb.ReturnOrderRequest{
		Order: 101,
	}
	_, err := client.ReturnOrder(ctx, req)
	if err != nil {
		log.Fatalf("could not return order: %v", err)
	}
	fmt.Println("Order returned")
}
