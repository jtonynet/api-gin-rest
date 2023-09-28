#!/bin/sh

ln -sf /usr/share/zoneinfo/America/Sao_Paulo /etc/localtime

if [ ! -d "bundle/bin" ]; then
    cd bundle

    echo "Downloading Gatling bundle..."
    wget https://repo1.maven.org/maven2/io/gatling/highcharts/gatling-charts-highcharts-bundle/3.9.5/gatling-charts-highcharts-bundle-3.9.5-bundle.zip

    echo "Unzip Gatling bundle..."
    unzip gatling-charts-highcharts-bundle-3.9.5-bundle.zip

    echo "Remove zip bundle..."
    rm -rf gatling-charts-highcharts-bundle-3.9.5-bundle.zip

    cd ..

    echo "Populate folder bundle..."
    mv bundle/gatling-charts-highcharts-bundle-3.9.5/* bundle

    echo "Remove original gatling folder..."
    chmod -R 777 bundle
    rm -rf bundle/gatling-charts-highcharts-bundle-3.9.5

    sleep 5
fi

echo "EXECUTE GatlingTest..."
$(pwd)/bundle/bin/gatling.sh -rm local -rd GatlingTest-$(exec date "+%m/%d/%Y-%H:%M:%S") -sf $(pwd)/APISimulation -rf $(pwd)/results APISimulation.Alunos

nc -l -p 8080
