package executor

import (
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

type CommandExecutorCli interface {
	GetCommands() []commands.Command
	Execute(cmdName string, args map[string]string) ([]models.Order, error)
}

type CommandExecutorGrpc interface {
	AcceptOrder(commands.AcceptOrder) ([]models.Order, error)
	AcceptReturn(commands.AcceptReturn) ([]models.Order, error)
	DeliverOrder(commands.DeliverOrder) ([]models.Order, error)
	GetOrders(commands.GetOrders) ([]models.Order, error)
	GetReturns(commands.GetReturns) ([]models.Order, error)
	ReturnOrder(commands.ReturnOrder) ([]models.Order, error)
}
