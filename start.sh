#!/bin/zsh

if [[ $# -ne 1 ]];then
  echo "Only support one parameter:(nodeId)"
  exit
else
  nodeId=$1
fi

go run main.go -me "${nodeId}" -role Follower -peers node1@localhost:7701,node2@localhost:7702,node3@localhost:7703