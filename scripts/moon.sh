#!/usr/bin/env bash
ssh -i ~/farza-key.pem ec2-user@18.224.114.44 "\
cd farza.party
git pull
sudo docker stop farza || true && sudo docker rm farza || true
sudo service docker start
sudo docker build -t nicememes .
sudo docker run -d --name=farza -p 80:80 nicememes
"