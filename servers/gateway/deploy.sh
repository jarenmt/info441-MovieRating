#!/bin/bash
. build.sh
docker push tilleyjaren/server
ssh ec2-user@ec2-3-138-193-149.us-east-2.compute.amazonaws.com < update.sh
