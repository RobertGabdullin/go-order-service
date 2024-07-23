package cli

import (
	"fmt"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/executor"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/logger"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/parser"
)

type CLI struct {
	parser   parser.Parser
	logger   logger.Logger
	executor executor.CommandExecutorCli
}

func NewCLI(executor executor.CommandExecutorCli, parser parser.Parser, logger logger.Logger) CLI {
	return CLI{
		parser:   parser,
		logger:   logger,
		executor: executor,
	}
}

func (c CLI) Help() {
	fmt.Println("Утилита для управления ПВЗ. Для аргументов команд можно использовать следующие форматы: -word=x --word=x -word x --word x. В примерах будет использован только формат -word=x. Список команд:")
	commands := c.executor.GetCommands()
	for i, elem := range commands {
		fmt.Printf("%d) Команда: %s\n   Описание: %s\n", i+1, elem.GetName(), elem.Description())
	}
}

func (c CLI) Run(input string) error {
	cmdName, mapArgs, err := c.parser.Parse(input)
	if err != nil {
		return err
	}

	event := logger.APIEvent{
		Timestamp:  time.Now(),
		MethodName: cmdName,
		RawRequest: input,
	}

	err = c.logger.Log(event)
	if err != nil {
		fmt.Printf("Failed to log event: %v", err)
	}

	_, err = c.executor.Execute(cmdName, mapArgs)
	return err
}
