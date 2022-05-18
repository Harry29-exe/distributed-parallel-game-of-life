#!/bin/bash

server="./server_main/server.out boardW=10000 boardH=10000 programIterations=10"
client1="./client/client.out"
client2="./client/client.out"

$server &

sleep 0.1s

$client1 &
$client2 &
