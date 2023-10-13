#!/bin/sh

if [ ! -e /entrypoint ]; then
    ln -sf /usr/share/zoneinfo/America/Sao_Paulo /etc/localtime
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

    echo "EXECUTE Gatling Test..."
    description=LoadTest::$API_NAME::v$API_TAG_VERSION::$(exec date "+%m/%d/%Y-%H:%M:%S")::America/Sao_Paulo
    sh $(pwd)/bundle/bin/gatling.sh -rm local -rd $description -sf $(pwd)/user-files/simulations/$API_NAME -rsf $(pwd)/user-files/resources/$API_NAME -rf $(pwd)/results/history

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

if [ "$1" = "clean-db" ]; then
    read -p "Do you really want to clean the database? This will reset your DB to zero records [y/n]: " answer
    if [ "$answer" = "y" ]; then
        export PGPASSWORD=$DATABASE_PASSWORD
        psql -h $DATABASE_HOST -U $DATABASE_USER -d $DATABASE_DB -p $DATABASE_PORT -c "TRUNCATE TABLE alunos RESTART IDENTITY;"
        unset PGPASSWORD
        echo "Database Cleaned"
    else
        echo "Database cleaning canceled."
    fi
fi


rm -rf ./results/latest/*
touch ./results/latest/.keep

new_latest=$(ls -td ./results/history/*/ | head -n 1)
cp -r $new_latest/* ./results/latest/

python3_pid=$(pgrep -f "python3 -m http.server $GATLING_TEST_PORT")
if [ ! -n "$python3_pid" ]; then
    echo "Run test result server"
    python3 -m http.server $GATLING_TEST_PORT --directory ./results/latest/
fi
