FROM quay.io/samsung_cnct/goglide:1.9.0

ENV PACKAGE_PATH "/go/src/github.com/samsung-cnct/crash-app-backend"
ENV     ELASTICSEARCH_TARGET_URL=http://elasticsearch:9200
ENV     RATE_LIMIT_PER_MINUTE=60

RUN apt-get -qq update && apt-get install -y -q build-essential && apt-get install file

WORKDIR ${PACKAGE_PATH}
COPY . ${PACKAGE_PATH}

RUN glide up
RUN make --no-builtin-rules --file make.golang build-app
RUN ls ./_containerize 
EXPOSE  8081
CMD  ./_containerize/crashbackend-linux serve --target $ELASTICSEARCH_TARGET_URL --ratelimit $RATE_LIMIT_PER_MINUTE
