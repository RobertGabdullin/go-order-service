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

	storageJSON := storage.NewStorage(fileName)
	parser := parser.MyParser{}
	cmd := cli.NewCLI(storageJSON, parser)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil || strings.TrimSpace(input) == "exit" {
			fmt.Println("End of program")
			os.Exit(0)
		} else {
			errRun := cmd.Run(input)
			if errRun != nil {
				fmt.Println(errRun)
			} else {
				fmt.Println("Success!")
			}
		}

	}
}
