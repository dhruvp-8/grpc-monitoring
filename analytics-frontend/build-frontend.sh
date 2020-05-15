#!/bin/bash

echo "creating the build file for the project"
yarn build
echo "moving the build file to ../src/static"
rm -rf ../src/public/build
mv build ../src/public

