FROM openjdk:8-jdk-alpine

WORKDIR /usr/src/app

COPY . .

RUN apk update && apk add --no-cache \
    python3 \
    tzdata \
    wget \
 && rm -rf /var/cache/apk/*

RUN chmod +x entrypoint.sh

ENTRYPOINT [ "./entrypoint.sh" ]
