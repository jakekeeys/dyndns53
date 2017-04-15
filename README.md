# dyndns53

A simple Go project to dynamically update and maintain a AWS route53 record, the value of which is based on the hosts current external IP.

## Building

To build the project run `go build .` followed by `docker build -t dydns53 .`

## Running

To run the project in docker use the following 
```
docker run -e POLL_INTERVAL_SECONDS='60' \
 -e IP_SERVICE_URL='http://checkip.amazonaws.com' \
 -e AWS_REGION='eu-west-1' \
 -e AWS_ACCESS_KEY='' \ 
 -e AWS_SECRET_KEY='' \
 -e AWS_HOSTED_ZONE_ID='' \
 -e AWS_RECORD_NAME='' \
 -e AWS_RECORD_TYPE='' dyndns53
``` 

To run the project locally use the following 
```
./dyndns53 --poll-interval-seconds=60 \
 --ip-service-url="http://checkip.amazonaws.com" \
 --aws-region="eu-west-1" \
 --aws-access-key="" \
 --aws-secret-key="" \
 --aws-hosted-zone-id="" \
 --aws-record-name="" \
 --aws-record-type=""
 ``` 

## Help

For help use the following `./dyndns53 help`