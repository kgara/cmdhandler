# cmdhandler-producer
Usage:
```
NAME:
   cmdhandler-producer - add commands from the config file to the ampq

USAGE:
   cmdhandler-producer [global options]

VERSION:
   1.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --ampq value      ampq uri string (default: "amqp://guest:guest@localhost:5672/")
   --queue value     ampq queue name (default: "job_queue")
   --scenario value  scenario.json file name (default: "/tmp/scenario01.json")
   --workers value   If we care about the order of the commands in the exact scenario - we might want to make one gorutine pool here.
                               Though again we should reduce the pool to a single gorutine on the consumer side as well.
                               Buy defaults we assume that we care for maintaining the order only in the orderedMap on the consumer. (default: 8)
   --help, -h        show help
   --version, -v     print the version
```