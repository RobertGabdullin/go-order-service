package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

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

	commandChan := make(chan string, 10)
	var wg sync.WaitGroup

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	numWorkers := 2
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(commandChan, &cmd)
		}()
	}

	go func() {
		<-sigs
		fmt.Println("\nReceived shutdown signal")
		close(commandChan)
	}()

	go func() {
		for {
			fmt.Print("> ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			input = strings.TrimSpace(input)
			if input == "exit" {
				fmt.Println("End of program")
				close(commandChan)
				return
			}
			commandChan <- input
		}
	}()

	wg.Wait()
}

func worker(commandChan <-chan string, cmd *cli.CLI) {
	for input := range commandChan {
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
