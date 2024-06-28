package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/cli"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/parser"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	connUrl := os.Getenv("DATABASE_URL")
	if connUrl == "" {
		fmt.Println("DATABASE_URL environment variable is not set")
		return
	}
	postgresStorage, err := storage.NewOrderStorage(connUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	wrapperStorage, err := storage.NewWrapperStorage(connUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	orderService := service.NewPostgresService(postgresStorage, wrapperStorage)

	parser := parser.ArgsParser{}
	cmd := cli.NewCLI(orderService, parser)

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
