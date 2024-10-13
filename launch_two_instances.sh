#!/bin/bash
# Script that launches two games at once (for multiplayer testing)

cleanup() {
    echo "Stopping processes..."
    kill $pid1 $pid2 2>/dev/null
    wait $pid1 2>/dev/null
    wait $pid2 2>/dev/null
    exit
}

# Trap for Ctrl+C to clean up processes
trap cleanup INT

# Start the processes in the background
process1_command="go run cmd/main.go --screen multiplayerJoin"
$process1_command 2>&1 &
pid1=$!

process2_command="go run cmd/main.go --screen multiplayerHost"
$process2_command 2>&1 &
pid2=$!

# Wait for either process to finish
wait -n $pid1 $pid2
