#!/bin/bash

originDir=$(pwd)

sdksDir=~/projects/go/src/github.com/nickolasgough/cloud-community-iam/sdks
sdkDir=${sdksDir}/typescript
cd ${sdkDir}

rm -rf ./src

echo "Generating TypeScript files..."
openapi-generator generate -g typescript-angular -i ../api.yaml -o ./src --additional-properties npmName=@cloud-community/iam,ngVersion=18.2.0

cd ./src
echo "Compiling TypeScript files..."
npm install
npm run build

echo "Publish the SDK..."
cd ./dist
npm publish --access public

cd ${originDir}
