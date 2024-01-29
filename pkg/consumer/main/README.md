# cmdhandler-consumer
Usage:
```
NAME:
cmdhandler-consumer - reads and executes commands from the ampq

USAGE:
cmdhandler-consumer [global options]

VERSION:
1.0

COMMANDS:
help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
--ampq value     ampq uri string (default: "amqp://guest:guest@localhost:5672/")
--queue value    ampq queue name (default: "job_queue")
--output value   processing output file name (default: "/tmp/consumer-output.txt")
--workers value  If we care about the order of the commands in the exact scenario - we might want to make one gorutine pool here.
And same on the producer side.
Buy defaults we assume that we care for maintaining the order only in the orderedMap on the consumer. (default: 8)
--help, -h       show help
--version, -v    print the version
```