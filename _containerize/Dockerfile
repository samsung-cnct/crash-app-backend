# Create a minimal golang app only container
FROM alpine:3.5
MAINTAINER Cindy ONeill <cubiedoo2@gmail.com>

ENV     ELASTICSEARCH_TARGET_URL=http://elasticsearch:9200
COPY    ./crashbackend-linux /crashbackend
EXPOSE  8081
CMD     /crashbackend serve --target $ELASTICSEARCH_TARGET_URL
