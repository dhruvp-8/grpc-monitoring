#!/bin/bash

echo "creating the build file for the project"
yarn build
echo "moving the build file to ../src/static"
mv build ../src/static
