FROM golang:latest

WORKDIR /usr/src/app

COPY . .

COPY entrypoint.sh ./ 

RUN chmod +x entrypoint.sh

ENTRYPOINT [ "./entrypoint.sh" ]