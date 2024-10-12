#!/bin/bash
# Script that launches two games at once (for multiplayer testing)

# Trap for Ctrl+C to clean up processes
trap cleanup INT

# Start the processes in the background
process1_command="go run cmd/main.go --screen multiplayerHost"
$process1_command &
pid1=$!

process2_command="go run cmd/main.go --screen multiplayerJoin"
$process2_command &
pid2=$!

# Wait for either process to finish
wait -n $pid1 $pid2

# If one process ends, clean up both
    echo "Stopping processes..."
    kill $pid1 $pid2 2>/dev/null
    wait $pid1 2>/dev/null
    wait $pid2 2>/dev/null
    exit
