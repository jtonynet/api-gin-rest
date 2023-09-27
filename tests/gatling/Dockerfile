FROM ubuntu:latest

WORKDIR /usr/src/app

RUN apt-get update && \
    apt-get install -y netcat wget openjdk-8-jdk scala unzip && \
    rm -rf /var/lib/apt/lists/*

COPY entrypoint.sh ./ 

RUN chmod +x entrypoint.sh

ENTRYPOINT [ "./entrypoint.sh" ]


# RUN ./bundle/bin/gatling.sh -sf $(pwd)/APISimulation -rf $(pwd)/results APISimulation.Alunos
