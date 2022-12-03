#!/bin/ash

redis-server & # redis
P1=$!
/app/bin/tenidditerApi & # api
P2=$!
wait $P1 $P2
