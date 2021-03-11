#!/bin/bash
docker rm -f test
docker pull tilleyjaren/a2clientfinal
export TLSCERT=/etc/letsencrypt/live/www.jaren441.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/www.jaren441.me/privkey.pem
docker run -d \
    --name test \
    -p 80:80 -p 443:443 \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    tilleyjaren/a2clientfinal
exit
