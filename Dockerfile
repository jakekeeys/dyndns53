FROM alpine

RUN apk add --no-cache ca-certificates

RUN mkdir /app

WORKDIR /app

ADD dyndns53 dyndns53

CMD ./dyndns53