#!/bin/sh

cd bundle

echo "Downloading Gatling bundle..."
wget https://repo1.maven.org/maven2/io/gatling/highcharts/gatling-charts-highcharts-bundle/3.9.5/gatling-charts-highcharts-bundle-3.9.5-bundle.zip

echo "Unzip Gatling bundle..."
unzip gatling-charts-highcharts-bundle-3.9.5-bundle.zip

echo "remove zip bundle..."
rm -rf gatling-charts-highcharts-bundle-3.9.5-bundle.zip

cd ..

echo "Popula o folder bundle..."
mv bundle/gatling-charts-highcharts-bundle-3.9.5/* bundle

echo "remove pasta original do gatling..."
chmod -R 777 bundle
rm -rf bundle/gatling-charts-highcharts-bundle-3.9.5

# ./bundle/bin/gatling.sh -sf $(pwd)/APISimulation -rf $(pwd)/results APISimulation.Alunos