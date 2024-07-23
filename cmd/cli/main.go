package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	_ "github.com/lib/pq"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/cli"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/config"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/event_broker"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/executor"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/logger"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/parser"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config file: %v\n", err)
		return
	}

	var kafkaClient *event_broker.KafkaClient
	if cfg.App.OutputMode == config.KafkaOutputMode {
		kafkaClient, err = event_broker.NewKafkaClient(cfg.Kafka.Brokers, nil)
		if err != nil {
			fmt.Printf("Error initializing Kafka: %v\n", err)
			return
		}

		defer kafkaClient.CloseProducer()
		go event_broker.StartConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	}

	connUrl := cfg.Database.URL
	if connUrl == "" {
		fmt.Println("Database URL is not set in the config file")
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
	executor := executor.NewOrderCommandExecutor(orderService)

	parser := parser.ArgsParser{}
	logger := logger.KafkaLogger{
		OutputMode:  cfg.App.OutputMode,
		KafkaTopic:  cfg.Kafka.Topic,
		KafkaClient: kafkaClient,
	}
	cmd := cli.NewCLI(executor, parser, logger)

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

	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			fmt.Print("> ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
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
