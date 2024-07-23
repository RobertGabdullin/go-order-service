package cli_grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/logger"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/parser"
	pb "gitlab.ozon.dev/r_gabdullin/homework-1/pb"
)

type CLI_gRPC struct {
	parser   parser.Parser
	client   pb.OrderServiceClient
	commands []commands.Command
}

func NewCLI(server pb.OrderServiceClient, parser parser.Parser, logger logger.Logger) CLI_gRPC {
	return CLI_gRPC{
		commands: []commands.Command{
			commands.AcceptOrder{},
			commands.AcceptReturn{},
			commands.DeliverOrder{},
			commands.GetOrders{},
			commands.GetReturns{},
			commands.ReturnOrder{},
		},
		parser: parser,
		client: server,
	}
}

func (c CLI_gRPC) Help() {
	fmt.Println("Утилита для управления ПВЗ. Для аргументов команд можно использовать следующие форматы: -word=x --word=x -word x --word x. В примерах будет использован только формат -word=x. Список команд:")
	for i, elem := range c.commands {
		fmt.Printf("%d) Команда: %s\n   Описание: %s\n", i+1, elem.GetName(), elem.Description())
	}
}

func (c CLI_gRPC) translate_command(cmd string, args map[string]string) ([]models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	switch cmd {
	case "acceptOrder":
		cmdArgs, err := commands.AcceptOrderAssignArgs(args)
		if err != nil {
			return nil, err
		}

		req := &pb.AcceptOrderRequest{
			User:      int32(cmdArgs.Order),
			Order:     int32(cmdArgs.Recipient),
			Weight:    int32(cmdArgs.Weight),
			BasePrice: int32(cmdArgs.BasePrice),
			Expire:    cmdArgs.Expire.String(),
			Wrapper:   cmdArgs.Wrapper,
		}

		_, err = c.client.AcceptOrder(ctx, req)
		if err != nil {
			return nil, err
		}
		return []models.Order{}, nil

	case "acceptReturn":
		cmdArgs, err := commands.AcceptReturnAssignArgs(args)
		if err != nil {
			return nil, err
		}

		req := &pb.AcceptReturnRequest{
			User:  int32(cmdArgs.User),
			Order: int32(cmdArgs.Order),
		}

		_, err = c.client.AcceptReturn(ctx, req)
		if err != nil {
			return nil, err
		}
		return []models.Order{}, nil

	case "deliverOrder":
		cmdArgs, err := commands.DeliverOrderAssignArgs(args)
		if err != nil {
			return nil, err
		}

		ords := make([]int32, 0)
		for _, i := range cmdArgs.Ords {
			ords = append(ords, int32(i))
		}

		req := &pb.DeliverOrderRequest{
			Orders: ords,
		}

		_, err = c.client.DeliverOrder(ctx, req)
		if err != nil {
			return nil, err
		}
		return []models.Order{}, nil

	case "getOrders":
		cmdArgs, err := commands.GetOrdersAssignArgs(args)
		if err != nil {
			return nil, err
		}

		req := &pb.GetOrdersRequest{
			User:  int32(cmdArgs.User),
			Count: int32(cmdArgs.Count),
		}

		resp, err := c.client.GetOrders(ctx, req)
		if err != nil {
			return nil, err
		}

		ords := make([]models.Order, 0)
		for _, i := range resp.Orders {
			expire, _ := time.Parse("2006-01-02T15", i.Expire)
			ords = append(ords, models.Order{
				Id:        int(i.Id),
				Recipient: int(i.Recipient),
				Expire:    expire,
				Status:    i.Status,
			})
		}
		return ords, nil

	case "getReturns":
		cmdArgs, err := commands.GetReturnsAssignArgs(args)
		if err != nil {
			return nil, err
		}

		req := &pb.GetReturnsRequest{
			Offset: int32(cmdArgs.Offset),
			Limit:  int32(cmdArgs.Limit),
		}

		resp, err := c.client.GetReturns(ctx, req)
		if err != nil {
			return nil, err
		}

		ords := make([]models.Order, 0)
		for _, i := range resp.Returns {
			expire, _ := time.Parse("2006-01-02T15", i.Expire)
			returnedAt, _ := time.Parse("2006-01-02T15", i.ReturnedAt)
			ords = append(ords, models.Order{
				Id:         int(i.Id),
				Recipient:  int(i.Recipient),
				Expire:     expire,
				ReturnedAt: returnedAt,
			})
		}

		return ords, nil

	case "returnOrder":
		cmdArgs, err := commands.ReturnOrderAssignArgs(args)
		if err != nil {
			return nil, err
		}

		req := &pb.ReturnOrderRequest{
			Order: int32(cmdArgs.Order),
		}

		_, err = c.client.ReturnOrder(ctx, req)
		if err != nil {
			return nil, err
		}
		return []models.Order{}, nil

	default:
		return nil, errors.New("unknown command")
	}
}

func (c CLI_gRPC) Run(input string) error {
	cmdName, mapArgs, err := c.parser.Parse(input)
	if err != nil {
		return err
	}

	_, err = c.translate_command(cmdName, mapArgs)
	return err
}
