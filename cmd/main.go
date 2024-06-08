package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/cli"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/parser"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

const (
	fileName = "orders.json"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	storageJSON, errCreate := storage.New(fileName)
	if errCreate != nil {
		fmt.Println(errCreate)
		return
	}

	parser := parser.ArgsParser{}
	cmd := cli.NewCLI(storageJSON, parser)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if err != nil || input == "exit" {
			fmt.Println("End of program")
			return
		}
		if input == "help" {
			cmd.Help()
			continue
		}
		errRun := cmd.Run(input)
		if errRun != nil {
			fmt.Println(errRun)
		} else {
			fmt.Println("Success!")
		}

	}
}
