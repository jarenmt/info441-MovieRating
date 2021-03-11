./build.sh
docker push ssunni12/ssunni12-web
ssh ec2-user@ec2-52-91-175-164.compute-1.amazonaws.com
sudo docker rm -f ssunni12-web
sudo docker pull ssunni12/ssunni12-web
export TLSCERT=/etc/letsencrypt/live/sanjayunni.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/sanjayunni.me/privkey.pem
sudo docker run -d -p 443:443 -p 80:80 -v /etc/letsencrypt:/etc/letsencrypt:ro -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY --name ssunni12-web ssunni12/ssunni12-web
