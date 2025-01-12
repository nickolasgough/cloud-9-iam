#!/bin/bash

sdksDir=~/projects/cloud-9/cloud-9-iam/sdks
sdkDir=${sdksDir}/typescript
cd ${sdkDir}

# Regenerate the SDK
rm -rf ./src

echo "Generating TypeScript files..."
openapi-generator generate -g typescript-angular -i ../api.yaml -o ./src --additional-properties npmName=@cloud-9/iam,ngVersion=18.2.0

cd ./src
echo "Compiling TypeScript files..."
npm install
npm run build

# Publish the SDK
echo "Publishing the SDK..."
cd ./dist
npm publish --access public
