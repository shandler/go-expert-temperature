#!/bin/bash


if [ -f .env ]; then
    set -o allexport
    source .env
    set +o allexport

    export BUILD_IMAGE=southamerica-east1-docker.pkg.dev/$PROJECT_ID/build/$IMAGE_PROD:latest

    docker build -t $BUILD_IMAGE -f Dockerfile.prod .
    docker tag $BUILD_IMAGE $BUILD_IMAGE
    docker push $BUILD_IMAGE

    IMAGE_TAG=southamerica-east1-docker.pkg.dev/$PROJECT_ID/build/$IMAGE_PROD
    echo "---------------------"
    echo $IMAGE_TAG
    echo "---------------------"

    TAGS=$(gcloud container images list-tags $IMAGE_TAG --format="get(tags)")
        echo "---------------------"
        echo $TAGS
        echo "---------------------"

        if [ -n "$TAGS" ]; then

        IFS=$'\n' read -rd '' -a TAG_ARRAY <<<"$TAGS"
        
        LAST_TAG=$(printf "%s\n" "${TAG_ARRAY[@]}" | sort -r | head -n 1)


        gcloud run deploy $IMAGE_PROD \
--image=southamerica-east1-docker.pkg.dev/$PROJECT_ID/build/$IMAGE_PROD:$LAST_TAG \
--allow-unauthenticated \
--port=8080 \
--service-account=$SERVICEACCOUNT \
--concurrency=10 \
--timeout=30 \
--memory=128Mi \
--set-env-vars=WEATHER_KEY=$WEATHER_KEY \
--cpu-boost \
--region=southamerica-east1 \
--project=$PROJECT_ID


    else
        echo "Erro: not found tag."
        exit 1
    fi
else
    echo "Erro: file .env not floud."
    exit 1
fi
