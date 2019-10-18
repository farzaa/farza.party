#!/usr/bin/env bash

# Build image locally, push it up to Dockerhub.
docker build -t party .
docker push farzatv/party

ssh -i ~/farza-key.pem ec2-user@18.224.114.44 "\
sudo docker stop farza || true && sudo docker rm farza || true
sudo service docker start
sudo docker pull farzatv/party
sudo docker run -d --name=farza -p 80:80 farzatv/party
"