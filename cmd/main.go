package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	_ "github.com/lib/pq"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/cli"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/parser"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
	"gitlab.ozon.dev/r_gabdullin/homework-1/kafka"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		URL string `yaml:"url"`
	} `yaml:"database"`
	Kafka struct {
		Brokers []string `yaml:"brokers"`
		Topic   string   `yaml:"topic"`
		GroupID string   `yaml:"group_id"`
	} `yaml:"kafka"`
	App struct {
		OutputMode string `yaml:"output_mode"`
	} `yaml:"app"`
}

func loadConfig() (*Config, error) {
	f, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	if cfg.App.OutputMode != "direct" && cfg.App.OutputMode != "kafka" {
		return nil, errors.New("Unknown output mode")
	}
	return &cfg, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config file:", err)
		return
	}

	if config.App.OutputMode == "kafka" {
		if err := kafka.InitKafka(config.Kafka.Brokers); err != nil {
			fmt.Println("Error initializing Kafka:", err)
			return
		}
		defer kafka.CloseProducer()
		go kafka.StartConsumer(config.Kafka.Brokers, config.Kafka.Topic)
	}

	connUrl := config.Database.URL
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

	parser := parser.ArgsParser{}
	cmd := cli.NewCLI(orderService, parser, config.App.OutputMode, config.Kafka.Topic)

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
