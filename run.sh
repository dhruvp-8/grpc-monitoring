#!/bin/bash
YELLOW='\033[1;33m'
CYAN='\033[1;36m'
ORANGE='\033[0;33m'
GREEN='\033[0;32m'
NC='\033[0m'

cd analytics-frontend

echo "${CYAN}Creating the build file for the front-end...${NC}\n"
yarn build

echo "\n\n${CYAN}Creating a public directory for storing build files (if not exists)...${NC}\n"
mkdir -p ../src/public

echo "${YELLOW}Removing the content inside the public directory (if already exists)...${NC}\n"
rm -r -f ../src/public/*

echo "${YELLOW}Moving the build files to public directory${NC}"
mv build/* ../src/public

rm -r -f build

echo "\n${GREEN}Instantiating the server${NC}"
cd ../src
go build
./grpc-analytics
