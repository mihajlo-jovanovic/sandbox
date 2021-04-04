#!/bin/zsh

echo "Building client image..."
docker build -f <(sed 's/api_impl/client/g' Dockerfile | sed 's/server/client/g' | sed '/EXPOSE/d') -t grpc_test_client .
echo "Done!"