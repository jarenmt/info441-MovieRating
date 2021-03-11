#!/bin/bash

./build.sh
docker push ssunni12/ssunni12
ssh ec2-user@ec2-54-173-220-226.compute-1.amazonaws.com < update.sh
