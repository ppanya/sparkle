#!/usr/bin/env bash
docker run -d --rm --name sparkle-pubsub -p 6690:8085 -e CLOUDSDK_CORE_PROJECT=test gcr.io/cloud-builders/gcloud beta emulators pubsub start --host-port=0.0.0.0:8085
docker run -d --rm --name sparkle-datastore -p 6691:5545 -e CLOUDSDK_CORE_PROJECT=test gcr.io/cloud-builders/gcloud beta emulators datastore start --host-port=0.0.0.0:5545

#docker run --name sparkle-mongo-primary \
#    --rm -p 6692:27017 \
#    -d \
#    -e MONGODB_ROOT_PASSWORD=root \
#    bitnami/mongodb:latest

#    -e MONGODB_REPLICA_SET_KEY=justprimary \
#    -e MONGODB_REPLICA_SET_MODE=primary \
#
#docker run --name sparkle-mongo-secondary \
#  --rm \
#  -d \
#  --link sparkle-mongo-primary:primary \
#  -e MONGODB_PRIMARY_ROOT_PASSWORD=root \
#  -e MONGODB_REPLICA_SET_KEY=justprimary \
#  -e MONGODB_REPLICA_SET_MODE=secondary \
#  -e MONGODB_PRIMARY_HOST=primary \
#  -e MONGODB_PRIMARY_PORT_NUMBER=27017 \
#  -p 6693:27017 \
#  bitnami/mongodb:latest
#
#docker run --name sparkle-mongo-secondary-2 \
#  --rm \
#  -d \
#  --link sparkle-mongo-primary:primary \
#  -e MONGODB_PRIMARY_ROOT_PASSWORD=root \
#  -e MONGODB_REPLICA_SET_KEY=justprimary \
#  -e MONGODB_REPLICA_SET_MODE=secondary \
#  -e MONGODB_PRIMARY_HOST=primary \
#  -e MONGODB_PRIMARY_PORT_NUMBER=27017 \
#  -p 6694:27017 \
#  bitnami/mongodb:latest
