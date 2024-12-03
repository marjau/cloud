#!/bin/bash

set -eux

export PATH=$PATH:${GCLOUD_DIR}/bin

PERSIST_DATA_FLAG="--no-store-on-disk"
if [ "$GCLOUD_STORE_ON_DISK" == true ]
then
    PERSIST_DATA_FLAG="--store-on-disk"    
fi

gcloud beta emulators datastore start \
    --project=$CLOUDSDK_CORE_PROJECT \
    --use-firestore-in-datastore-mode \
    --host-port=$DS_HOST:8081 \
    $PERSIST_DATA_FLAG
