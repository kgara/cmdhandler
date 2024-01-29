package main

import (
	"context"
	"encoding/json"
	"github.com/kgara/cmdhandler/pkg/consumer"
	"github.com/kgara/cmdhandler/pkg/shared"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

//const numWorkers = 2
//const ampqQueueName = "job_queue"
//const ampqAddr = "amqp://guest:guest@localhost:5672/"
//const processingFilename = "/tmp/consumer-output.txt"

type ConsumerConfig struct {
	ampqUri            string
	ampqQueueName      string
	processingFileName string
	numWorkers         int
}

func main() {
	config := &ConsumerConfig{}

	app := &cli.App{
		Name:      "cmdhandler-consumer",
		Usage:     "reads and executes commands from the ampq",
		Version:   "1.0",
		UsageText: "cmdhandler-consumer [global options]",
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
				Name:        "output",
				Value:       "/tmp/consumer-output.txt",
				Usage:       "processing output file name",
				Destination: &config.processingFileName,
			},
			&cli.IntFlag{
				Name:  "workers",
				Value: 8,
				Usage: `If we care about the order of the commands in the exact scenario - we might want to make one gorutine pool here.
						And same on the producer side.
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

func execute(config *ConsumerConfig) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	queue := shared.NewClient(config.ampqQueueName, config.ampqUri, logger)

	// Give the connection sometime to set up
	for !queue.IsReady {
		<-time.After(time.Second)
	}

	fileWriter := consumer.NewFileWriter(config.processingFileName, logger)
	fileWriter.Start()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	defer fileWriter.Close()
	defer logger.Println("Shutting down...")

	deliveries, err := queue.Consume()
	if err != nil {
		logger.Printf("Could not start consuming: %s\n", err)
		return
	}

	// This channel will receive a notification when a channel closed event
	// happens. This must be different from Client.notifyChanClose because the
	// library sends only one notification and Client.notifyChanClose already has
	// a receiver in handleReconnect().
	// Recommended to make it buffered to avoid deadlocks
	chClosedCh := make(chan *amqp.Error, 1)
	queue.Channel.NotifyClose(chClosedCh)

	// Create an OrderedMapV1 instance
	orderedMap := consumer.NewOrderedMap(fileWriter)

	// Start worker pool
	startWorkers(config, deliveries, orderedMap, logger)
	// Handle meta-situations
	for {
		select {
		case <-ctx.Done():
			_ = queue.Close()
			return

		case amqErr := <-chClosedCh:
			// This case handles the event of closed channel e.g. abnormal shutdown
			logger.Printf("AMQP Channel was closed due to: %s. Reinitializing consuming...\n", amqErr)

			deliveries, err = queue.Consume()
			if err != nil {
				// If the AMQP channel is not ready, it will continue the loop. Next
				// iteration will enter this case because chClosedCh is closed by the
				// library
				logger.Println("Error trying to consume, will try again")
				<-time.After(time.Second)
				continue
			}
			startWorkers(config, deliveries, orderedMap, logger)

			// Re-set channel to receive notifications
			// The library closes this channel after abnormal shutdown
			chClosedCh = make(chan *amqp.Error, 1)
			queue.Channel.NotifyClose(chClosedCh)
			logger.Printf("Consuming again!\n")
		}
	}
}

func processDelivery(delivery amqp.Delivery, orderedMap consumer.OrderedMap, logger *log.Logger) {
	logger.Printf("Received message: %s\n", delivery.Body)
	command := &shared.Command{}
	err := json.Unmarshal(delivery.Body, command)
	if err != nil {
		logger.Printf("Error decoding JSON: %s\n", err)
		err := delivery.Nack(false, false)
		if err != nil {
			logger.Printf("Error negatively acknowledging message: %s\n", err)
		}
		return
	}
	//logger.Printf("Received command: %s\n", *command)
	orderedMap.ExecuteCommand(command)
	if err := delivery.Ack(false); err != nil {
		logger.Printf("Error acknowledging message: %s\n", err)
	}
}

func startWorkers(config *ConsumerConfig, deliveries <-chan amqp.Delivery, orderedMap *consumer.OrderedMapImpl, logger *log.Logger) {
	for i := 0; i < config.numWorkers; i++ {
		go worker(i, deliveries, orderedMap, logger)
	}
}

func worker(id int, deliveries <-chan amqp.Delivery, orderedMap consumer.OrderedMap, logger *log.Logger) {
	for delivery := range deliveries {
		logger.Printf("Worker %d: received task\n", id)
		processDelivery(delivery, orderedMap, logger)
	}
	logger.Printf("Worker %d: finished\n", id)
}
