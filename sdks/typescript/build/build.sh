#!/bin/bash

originDir=$(pwd)

sdksDir=~/projects/go/src/github.com/nickolasgough/cloud-9-iam/sdks
sdkDir=${sdksDir}/typescript
cd ${sdkDir}

rm -rf ./src

echo "Generating TypeScript files..."
openapi-generator generate -g typescript-angular -i ../api.yaml -o ./src --additional-properties npmName=@cloud-9/iam,ngVersion=18.2.0

cd ./src
echo "Compiling TypeScript files..."
npm install
npm run build

echo "Publishing the SDK..."
cd ./dist
npm publish --access public

cd ${originDir}
