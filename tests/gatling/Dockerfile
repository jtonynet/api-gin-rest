FROM ubuntu:latest

WORKDIR /usr/src/app

RUN apt-get update && \
    apt-get install -y netcat wget openjdk-8-jdk scala unzip && \
    rm -rf /var/lib/apt/lists/*


RUN echo "Downloading Gatling bundle..."
RUN wget https://repo1.maven.org/maven2/io/gatling/highcharts/gatling-charts-highcharts-bundle/3.9.5/gatling-charts-highcharts-bundle-3.9.5-bundle.zip

RUN echo "Unzip Gatling bundle..."
RUN unzip gatling-charts-highcharts-bundle-3.9.5-bundle.zip

RUN echo "Rename Gatling bundle..."
RUN mv gatling-charts-highcharts-bundle-3.9.5 bundle

# RUN ./bundle/bin/gatling.sh -sf $(pwd)/APISimulation -rf $(pwd)/results APISimulation.Alunos
