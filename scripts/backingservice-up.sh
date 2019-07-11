#!/usr/bin/env bash
docker run -d --rm --name sparkle-pubsub -p 6690:8085 -e CLOUDSDK_CORE_PROJECT=test gcr.io/cloud-builders/gcloud beta emulators pubsub start --host-port=0.0.0.0:8085
docker run -d --rm --name sparkle-datastore -p 6691:5545 -e CLOUDSDK_CORE_PROJECT=test gcr.io/cloud-builders/gcloud beta emulators datastore start --host-port=0.0.0.0:5545

docker run --name sparkle-mongo-primary \
    --rm -p 6692:27017 \
    -d \
    -e MONGODB_ROOT_PASSWORD=root \
    bitnami/mongodb:latest