#!/bin/bash

docker rm -f ssunni12
docker rm -f mysqlServer
docker rm -f redisServer
docker network rm 441network
docker network create 441network
sudo docker pull ssunni12/ssunni12
export TLSCERT=/etc/letsencrypt/live/api.sanjayunni.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.sanjayunni.me/privkey.pem
export SESSIONKEY=sessionkey
export DB_NAME=441sqlserver
export MYSQL_ROOT_PASSWORD=testPassword
export ADDR=:443
sudo docker run -d -p 6379:6379 --network 441network --name redisServer redis
export REDISADDR=redisServer:6379
sudo docker run -d --name mysqlServer  --network 441network -p 3306:3306 -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD -e MYSQL_DATABASE=$DB_NAME ssunni12/ssunni12sql
export DSN=root:$MYSQL_ROOT_PASSWORD@tcp\(mysqlServer:3306\)/$DB_NAME
sudo docker run -d -p 443:443 -e ADDR=$ADDR --name ssunni12 -v /etc/letsencrypt:/etc/letsencrypt:ro -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY -e DSN=$DSN -e SESSIONKEY=$SESSIONKEY --network 441network ssunni12/ssunni12
