#!/usr/bin/env bash

CUR_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

cd $CUR_DIR/../ 

tmux new-session -s chiro -d

tmux send "/usr/bin/go run ./p2sub --key-file ./json/node1.json --bind-port 4433 --bind-host 0.0.0.0" C-m
tmux rename-window "Program 1"
sleep 5
tmux split-window
tmux send "/usr/bin/go run ./p2sub --key-file ./json/node2.json --bind-port 4434 --bind-host 0.0.0.0" C-m
tmux rename-window "Program 2"
sleep 5
tmux split-window
tmux send "/usr/bin/go run ./p2sub --key-file ./json/node3.json --bind-port 4435 --bind-host 0.0.0.0" C-m
tmux rename-window "Program 3"
tmux attach