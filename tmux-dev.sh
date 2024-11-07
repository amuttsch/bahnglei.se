#!/bin/zsh

SESH="bahngleise"

tmux has-session -t $SESH 2>/dev/null

if [ $? != 0 ]; then
  tmux new-session -d -s $SESH -n "nvim"
  tmux send-keys -t $SESH:nvim "vim ." C-m

  tmux new-window -t $SESH -n "air"
  tmux send-keys -t $SESH:air "air serve" C-m

  tmux new-window -t $SESH -n "import"

  tmux new-window -t $SESH -n "osm"
  tmux send-keys -t $SESH:osm "cd osm" C-m
  tmux send-keys -t $SESH:osm "python3 -m http.server 8100" C-m

  tmux new-window -t $SESH -n "pg-fly"
  tmux send-keys -t $SESH:pg-fly "fly proxy 5433 -a bahngleise-db" C-m
fi

tmux attach-session -t $SESH

