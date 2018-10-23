#!/bin/sh

echo "=================================="
echo "building main application"
echo "=================================="
npm install
npm install semantic-ui --save
npm rebuild
npm run build
ls -al
rm -rf /opt/build/*
mv dist/* /opt/build/
rm -rf dist/
