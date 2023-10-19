FROM openjdk:17-alpine

WORKDIR /usr/src/app

COPY . .

RUN apk update && apk add --no-cache \
    postgresql-client \
    python3 \
    tzdata \
    wget \
 && rm -rf /var/cache/apk/*

RUN chmod +x entrypoint.sh

ENTRYPOINT [ "./entrypoint.sh" ]
