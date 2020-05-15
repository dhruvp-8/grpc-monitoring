#!/bin/bash

echo "creating the build file for the project"
yarn build
echo "moving the build file to ../src/public"
rm -rf ../src/public/
mv build ../src/public

