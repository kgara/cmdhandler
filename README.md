# cmdhandler

Simple AMPQ Consumer-Producer example with OrderedMap backend storage.

AMPQ interaction code is based on the reference guidence from
https://pkg.go.dev/github.com/rabbitmq/amqp091-go

Tested with RabbitMQ 3.9.13 listening on http://localhost:5672/

For emulating the "slow io operations" and parallelism demonstration the writer.Write operation may be "intentionally" delayed, so we can demonstrate our 
"delivery processing" in a thread pool.

# Components

* [cmdhandler-consumer](pkg/consumer/main) - reads and executes commands from the ampq
* [cmdhandler-producer](pkg/producer/main) - add commands from the config file to the ampq


# Building

    $ git clone https://github.com/kgara/cmdhandler.git
    $ cd cmdhandler && make

# Test run example scenario

    # Install and start RabbitMQ 
    # (out of the box http://localhost:5672/ with guest:guest)
    # Start consumer:
    $ ./build/cmdhandler-consumer
    # May start to observe the default output file in another console:
    $ touch /tmp/consumer-output.txt && tail -f /tmp/consumer-output.txt
    # Execute producer as many times as required:
    $ ./build/cmdhandler-producer --scenario=./pkg/producer/examples/scenario01.json
    # Also can run multiple producers in parallel, e.g. like that:
    $ ./scripts/multiproducers.sh
