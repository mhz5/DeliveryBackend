#!/bin/bash

cd ..

# Bundle all .proto files into a single compiled .pb file.
OUT_PB="./out.pb"
protoc \
        --include_imports \
        --include_source_info \
        --descriptor_set_out ${OUT_PB} \
        proto/*/*.proto

# Deploy the Cloud Endpoints service
gcloud endpoints services deploy ${OUT_PB} ./scripts/api_config.yaml --verbosity=debug

rm ${OUT_PB}
