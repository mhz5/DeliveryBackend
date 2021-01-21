#!/usr/bin/env bash

INSTANCE_NAME="grpc-host"
ON_VM_DEPLOY_SCRIPT_NAME="ON_VM_deploy_food_grpc_server.sh"
FDA_GRPC_IMAGE="fda-grpc-server"

# Build the project
./generate_buildfiles_and_build.sh

if [[ $? -ne 0 ]]; then
    exit 1
fi

cd ..

# Build docker image of server, then tag and push to google Cloud Registry.
# Don't forget the _cgo suffix to compile with cgo!
bazel build //server:docker_server.tar --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64_cgo
if [[ $? -ne 0 ]]; then
    exit 1
fi
docker load -i bazel-bin/server/docker_server.tar
docker tag bazel/server:docker_server gcr.io/food-prod/${FDA_GRPC_IMAGE}:latest
docker push gcr.io/food-prod/${FDA_GRPC_IMAGE}:latest

# Copy the deployment script to the VM.
gcloud compute scp --zone=us-central1-a "scripts/${ON_VM_DEPLOY_SCRIPT_NAME}" ${INSTANCE_NAME}:./scripts

# Deploy server on VM.
gcloud compute ssh ${INSTANCE_NAME} --zone=us-central1-a --command="bash ./scripts/${ON_VM_DEPLOY_SCRIPT_NAME}"
