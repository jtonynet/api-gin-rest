#!/bin/sh

ln -sf /usr/share/zoneinfo/America/Sao_Paulo /etc/localtime

if [ ! -e /entrypoint ]; then
    ln -s /usr/src/app/entrypoint.sh /entrypoint
fi

if [ "$1" = "run-test" ]; then

    if [ ! -d "bundle/bin" ]; then

        cd bundle

        echo "Downloading Gatling bundle..."
        wget  https://repo1.maven.org/maven2/io/gatling/highcharts/gatling-charts-highcharts-bundle/3.9.5/gatling-charts-highcharts-bundle-3.9.5-bundle.zip

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

    fi

    echo "EXECUTE GatlingTest..."
    description=LoadTest::$API_NAME::v$API_TAG_VERSION::$(exec date "+%m/%d/%Y-%H:%M:%S")::America/Sao_Paulo
    $(pwd)/bundle/bin/gatling.sh -rm local -rd $description -sf $(pwd)/APISimulation -rf $(pwd)/results/history APISimulation.Alunos
    echo "Verify Test Gatling Results folder for all tests"
fi

if [ "$1" = "clean-test" ]; then
    rm -rf ./bundle/*
    touch ./bundle/.keep

    directory="./results/history/"
    keep_folder="default"
    for item in "$directory"/*; do
        if [ -d "$item" ] && [ "$(basename "$item")" != "$keep_folder" ]; then
            rm -rf "$item"
        fi
    done
fi

echo "Remove last load test data"
rm -rf ./results/latest/*
touch ./results/latest/.keep

echo "Add New load test data"
new_latest=$(ls -td ./results/history/*/ | head -n 1)
cp -r $new_latest/* ./results/latest/

python3_pid=$(pgrep -f "python3 -m http.server $TEST_GATLING_PORT")
if [ ! -n "$python3_pid" ]; then
    echo "Run test last result server"
    python3 -m http.server $TEST_GATLING_PORT --directory ./results/latest/
fi


