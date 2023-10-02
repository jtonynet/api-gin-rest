FROM golang:1.21.1

WORKDIR /usr/src/app

COPY . . 

RUN chmod +x entrypoint.sh

ENTRYPOINT [ "./entrypoint.sh" ]