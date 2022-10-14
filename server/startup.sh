#!/bin/ash

/app/bin/tenidditerApi & # api
P1=$!
redis-server & # redis
P2=$!
wait $P1 $P2
