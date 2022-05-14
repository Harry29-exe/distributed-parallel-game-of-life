#!/bin/bash

server="./server_main/server.out boardW=120 boardH=120 programIterations=100 delay=true delayTime=10"
client1="./client/client.out"
client2="./client/client.out"

$server &

sleep 0.1s

$client1 &
$client2 &
