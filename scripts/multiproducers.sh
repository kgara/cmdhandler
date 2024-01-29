#!/bin/bash

# Number of instances to run in parallel
num_instances=10

# Define the command to run
command="./build/cmdhandler-producer --scenario=./pkg/producer/examples/scenario01.json --workers=4"

# Run instances in parallel
for ((i=1; i<=$num_instances; i++)); do
    # Run command in the background
    $command &
done

# Wait for all instances to finish
wait