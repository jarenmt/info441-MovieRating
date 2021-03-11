#!/bin/bash
. build.sh
docker push tilleyjaren/a2clientfinal
ssh ec2-user@ec2-18-222-64-36.us-east-2.compute.amazonaws.com < update.sh