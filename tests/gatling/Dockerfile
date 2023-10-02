FROM openjdk:8

WORKDIR /usr/src/app

COPY . .

RUN apt-get update && apt-get install -y \
    python3 \
    tzdata \
    wget \
 && apt-get clean

RUN chmod +x entrypoint.sh

ENTRYPOINT [ "./entrypoint.sh" ]
