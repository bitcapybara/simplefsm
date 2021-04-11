#!/bin/zsh

if [[ $# -ne 1 ]];then
  echo "Only support one parameter:(nodeId)"
  exit
else
  nodeId=$1
fi

go run server.go -me "${nodeId}" -role Learner -peers node1@localhost:7701,node2@localhost:7702,node3@localhost:7703,node4@localhost:7704,node5@localhost:7705