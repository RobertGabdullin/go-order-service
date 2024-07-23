package executor

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
	"gitlab.ozon.dev/r_gabdullin/homework-1/pkg/hash"
)

type orderCommandExecutor struct {
	service  service.StorageService
	mu       *sync.Mutex
	commands []commands.Command
}

func NewOrderCommandExecutor(service service.StorageService) orderCommandExecutor {
	return orderCommandExecutor{
		service: service,
		commands: []commands.Command{
			commands.AcceptOrder{},
			commands.AcceptReturn{},
			commands.DeliverOrder{},
			commands.GetOrders{},
			commands.GetReturns{},
			commands.ReturnOrder{},
		},
		mu: &sync.Mutex{},
	}
}

func (e orderCommandExecutor) GetCommands() []commands.Command {
	return e.commands
}

func (e orderCommandExecutor) Execute(cmd string, args map[string]string) ([]models.Order, error) {
	switch cmd {
	case "acceptOrder":
		cmdArgs, err := commands.AcceptOrderAssignArgs(args)
		if err != nil {
			return nil, err
		}

		ords, err := e.AcceptOrder(cmdArgs)
		if err != nil {
			return nil, err
		}
		return ords, nil
	case "acceptReturn":
		cmdArgs, err := commands.AcceptReturnAssignArgs(args)
		if err != nil {
			return nil, err
		}

		ords, err := e.AcceptReturn(cmdArgs)
		if err != nil {
			return nil, err
		}
		return ords, nil
	case "deliverOrder":
		cmdArgs, err := commands.DeliverOrderAssignArgs(args)
		if err != nil {
			return nil, err
		}

		ords, err := e.DeliverOrder(cmdArgs)
		if err != nil {
			return nil, err
		}
		return ords, nil

	case "getOrders":
		cmdArgs, err := commands.GetOrdersAssignArgs(args)
		if err != nil {
			return nil, err
		}

		ords, err := e.GetOrders(cmdArgs)
		if err != nil {
			return nil, err
		}
		return ords, nil

	case "getReturns":
		cmdArgs, err := commands.GetReturnsAssignArgs(args)
		if err != nil {
			return nil, err
		}

		ords, err := e.GetReturns(cmdArgs)
		if err != nil {
			return nil, err
		}
		return ords, nil

	case "returnOrder":
		cmdArgs, err := commands.ReturnOrderAssignArgs(args)
		if err != nil {
			return nil, err
		}

		ords, err := e.ReturnOrder(cmdArgs)
		if err != nil {
			return nil, err
		}
		return ords, nil

	default:
		return nil, errors.New("unknown command")
	}
}

func (e orderCommandExecutor) AcceptOrder(args commands.AcceptOrder) ([]models.Order, error) {
	if args.Expire.Before(time.Now()) {
		return nil, errors.New("storage time is out")
	}

	wrapper, err := e.service.GetWrapper(args.Wrapper)
	if err != nil {
		return nil, err
	}

	if wrapper.MaxWeight.Valid && args.Weight > int(wrapper.MaxWeight.Int64) {
		return nil, errors.New("order weight exceeds the maximum limit for the chosen wrapper")
	}

	hash := hash.GenerateHash()

	e.mu.Lock()
	defer e.mu.Unlock()

	err = e.service.AddOrder(models.NewOrder(args.Order, args.Recipient, args.Expire, "alive", hash, args.Weight, args.BasePrice, wrapper.Type))
	if err != nil {
		return nil, err
	}

	return []models.Order{}, nil
}

func (e orderCommandExecutor) ReturnOrder(args commands.ReturnOrder) ([]models.Order, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	temp := []int{args.Order}
	ords, err := e.service.FindOrders(temp)

	if err != nil {
		return nil, err
	}

	if len(ords) == 0 {
		return nil, errors.New("such order does not exist")
	}

	if ords[0].Status != "alive" && ords[0].Status != "returned" {
		return nil, errors.New("order is not at storage")
	}

	if ords[0].Expire.After(time.Now()) {
		return nil, errors.New("order should be out of storage limit")
	}

	err = e.service.DeleteOrder(args.Order)
	if err != nil {
		return nil, err
	}

	return []models.Order{}, nil
}

func (e orderCommandExecutor) AcceptReturn(args commands.AcceptReturn) ([]models.Order, error) {
	hash := hash.GenerateHash()

	e.mu.Lock()
	defer e.mu.Unlock()

	temp := []int{args.Order}
	ords, err := e.service.FindOrders(temp)

	if err != nil {
		return nil, err
	}

	for _, elem := range ords {
		if elem.Id != args.Order {
			continue
		}
		if elem.Status != "delivered" {
			return nil, errors.New("such an order has never been issued")
		}
		if elem.DeliveredAt.AddDate(0, 0, 2).Before(time.Now()) {
			return nil, errors.New("the order can only be returned within two days after issue")
		}
		err = e.service.ChangeStatus(args.Order, "returned", hash)
		if err != nil {
			return nil, err
		}

		return []models.Order{}, nil
	}

	return nil, errors.New("order with such ids does not exist")
}

func (e orderCommandExecutor) DeliverOrder(args commands.DeliverOrder) ([]models.Order, error) {
	hash := hash.GenerateHash()

	e.mu.Lock()
	defer e.mu.Unlock()

	ords, err := e.service.FindOrders(args.Ords)
	if err != nil {
		return nil, err
	}

	temp := make(map[int]bool)
	for _, elem := range ords {
		temp[elem.Recipient] = true
	}
	if len(temp) > 1 {
		return nil, errors.New("list of orders should belong only to one person")
	}

	for _, elem := range ords {
		if elem.Status != "alive" {
			return nil, errors.New("some orders are not available")
		}
		if elem.Expire.Before(time.Now()) {
			return nil, errors.New("some orders is out of storage limit date")
		}
		tempErr := e.service.ChangeStatus(elem.Id, "delivered", hash)
		if tempErr != nil {
			return nil, tempErr
		}
	}

	return []models.Order{}, nil
}

func (e orderCommandExecutor) GetOrders(args commands.GetOrders) ([]models.Order, error) {
	e.mu.Lock()
	ords, err := e.service.ListOrders(args.User)
	e.mu.Unlock()

	if err != nil {
		return nil, err
	}

	cnt := 1

	ans := make([]models.Order, 0)
	for i := len(ords) - 1; i >= 0 && (args.Count == -1 || args.Count >= cnt); i-- {
		if ords[i].Status == "alive" {
			fmt.Printf("%d) orderID = %d recipientID = %d storedUntil = %s\n", cnt, ords[i].Id, ords[i].Recipient, ords[i].Expire)
			cnt++
			ans = append(ans, ords[i])
		}
	}

	return ans, nil
}

func (e orderCommandExecutor) GetReturns(args commands.GetReturns) ([]models.Order, error) {
	e.mu.Lock()
	ords, err := e.service.GetReturns(args.Offset, args.Limit)
	e.mu.Unlock()

	if err != nil {
		return nil, err
	}

	for i := range ords {
		fmt.Printf("%d) orderID = %d recipientID = %d storedUntil = %s acceptedAt = %s\n", i+1, ords[i].Id, ords[i].Recipient, ords[i].Expire, ords[i].ReturnedAt)
	}

	return ords, nil
}
