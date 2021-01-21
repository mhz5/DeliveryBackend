#!/usr/bin/env bash

# This script belongs at ~/scripts on the Compute engine VM that hosts the gRPC server.

GOOGLE_CLOUD_PROJECT=$(curl -s "http://metadata.google.internal/computeMetadata/v1/project/project-id" -H "Metadata-Flavor: Google")
FDA_ENDPOINTS_NAME="fda.endpoints.${GOOGLE_CLOUD_PROJECT}.cloud.goog"
FDA_GRPC_NAME="fda-grpc-server"
GRPC_SERVER_IMAGE_NAME="gcr.io/${GOOGLE_CLOUD_PROJECT}/${FDA_GRPC_NAME}:latest"
ENDPOINTS_RUNTIME_IMAGE_NAME="gcr.io/endpoints-release/endpoints-runtime:2"
CORS_ALLOW_HEADERS=\
"Access-Control-Allow-Origin,\
X-User-Agent,\
X-Grpc-Web,\
Content-Type,\
X-Api-Key"

# These two steps only need to be done once to configure the server.
docker-credential-gcr configure-docker
docker network create --driver bridge esp_net

# Stop and remove all existing Docker containers.
docker stop $(docker ps -a -q); docker rm $(docker ps -a -q)

# Pull the newest versions of the relevant Docker images.
docker pull ${GRPC_SERVER_IMAGE_NAME}
docker pull ${ENDPOINTS_RUNTIME_IMAGE_NAME}

# Run the GRPC server.
docker run --detach --net=esp_net \
    --name=${FDA_GRPC_NAME} \
    ${GRPC_SERVER_IMAGE_NAME}

# Run the esp2 proxy.
# Startup options for esp2: https://cloud.google.com/endpoints/docs/openapi/specify-esp-v2-startup-options
docker run \
    --detach \
    --name=esp \
    --publish=443:9000 \
    --net=esp_net \
    ${ENDPOINTS_RUNTIME_IMAGE_NAME} \
    --cors_preset=basic \
    --cors_allow_headers=${CORS_ALLOW_HEADERS} \
    --service=${FDA_ENDPOINTS_NAME} \
    --rollout_strategy=managed \
    --backend=grpc://${FDA_GRPC_NAME}:50051 \
    --generate_self_signed_cert \
    -z healthz \
    --listener_port=9000

#    --volume=/etc/esp/ssl:/etc/esp/ssl \
