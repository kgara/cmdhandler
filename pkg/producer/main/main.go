package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/kgara/cmdhandler/pkg/shared"
	"github.com/urfave/cli/v2"
)

type ProducerConfig struct {
	ampqUri          string
	ampqQueueName    string
	scenarioFileName string
	numWorkers       int
}

func main() {
	config := &ProducerConfig{}

	app := &cli.App{
		Name:      "cmdhandler-producer",
		Usage:     "add commands from the config file to the ampq",
		Version:   "1.0",
		UsageText: "cmdhandler-producer [global options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "ampq",
				Value:       "amqp://guest:guest@localhost:5672/",
				Usage:       "ampq uri string",
				Destination: &config.ampqUri,
			},
			&cli.StringFlag{
				Name:        "queue",
				Value:       "job_queue",
				Usage:       "ampq queue name",
				Destination: &config.ampqQueueName,
			},
			&cli.StringFlag{
				Name:        "scenario",
				Value:       "/tmp/scenario01.json",
				Usage:       "scenario.json file name",
				Destination: &config.scenarioFileName,
			},
			&cli.IntFlag{
				Name:  "workers",
				Value: 8,
				Usage: `If we care about the order of the commands in the exact scenario - we might want to make one gorutine pool here.
						Though again we should reduce the pool to a single gorutine on the consumer side as well.
						Buy defaults we assume that we care for maintaining the order only in the orderedMap on the consumer.`,
				Destination: &config.numWorkers,
			},
		},
		Action: func(cCtx *cli.Context) error {
			execute(config)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
func execute(config *ProducerConfig) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	queue := shared.NewClient(config.ampqQueueName, config.ampqUri, logger)

	// Give the connection sometime to set up
	for !queue.IsReady {
		<-time.After(time.Second)
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute*10))
	defer cancel()
	defer logger.Println("Shutting down...")

	commandsChannel := make(chan shared.Command)
	var wg sync.WaitGroup
	// Start submitters pool
	for i := 0; i < config.numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for command := range commandsChannel {
				logger.Printf("Worker %d: received task\n", id)
				commandJson, err := json.Marshal(command)
				if err != nil {
					logger.Printf("Error encoding JSON: %s\n", err)
				}
				if err := queue.Push(commandJson); err != nil {
					logger.Printf("Push failed: %s\n", err)
				} else {
					logger.Println("Push succeeded!")
				}
			}
			logger.Printf("Worker %d: ended\n", id)
		}(i)
	}

	repeatableCommands, err := parseConfiguration(config.scenarioFileName, logger)
	if err != nil {
		return
	}

	for _, repeatableCommand := range repeatableCommands {
		// Add the Command to the channel 'Times' number of Times
		for i := 0; i < repeatableCommand.Times; i++ {
			commandsChannel <- repeatableCommand.Command
		}
	}
	close(commandsChannel)

	go func() {
		wg.Wait()
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			_ = queue.Close()
			return
		}
	}

}

func parseConfiguration(configurationFilename string, logger *log.Logger) (repeatableCommands []RepeatableCommand, err error) {
	file, err := os.Open(configurationFilename)
	if err != nil {
		logger.Println("Error opening file:", err)
		return nil, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		logger.Println("Error reading file:", err)
		return nil, err
	}

	err = json.Unmarshal(content, &repeatableCommands)
	if err != nil {
		logger.Println("Error unmarshalling JSON:", err)
		return nil, err
	}
	err = file.Close()
	if err != nil {
		logger.Println("Error closing the file:", err)
		return nil, err
	}

	return repeatableCommands, nil
}

type RepeatableCommand struct {
	Command shared.Command
	Times   int
}
