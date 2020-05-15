#!/bin/bash

echo "creating the build file for the project"
yarn build
echo "instantiating the server"
mkdir -p ../src/public
rm -r -f ../src/public/*
mv build/* ../src/public
rm -r -f build
cd ../src
go build
./grpc-analytics
