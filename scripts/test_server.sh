#!/usr/bin/env bash

HOST_IP="35.244.146.186:443"
cd ..
bazel run //client -- \
    --addr="${HOST_IP}" \
    --api-key=AIzaSyC4_1PlUx483TfgsA-zhKdPH3uop_iwwP8 \
    "test message"

