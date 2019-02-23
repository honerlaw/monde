#!/bin/sh

eval $(aws ecr get-login --no-include-email --region=us-east-1)

URL=106480132517.dkr.ecr.us-east-1.amazonaws.com/monde:$TAG

docker build -t $URL .
docker push $URL