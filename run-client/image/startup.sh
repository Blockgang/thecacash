#!/bin/sh

chown -R root:root /root

rm -rf app/
mkdir -p app/
cp build/* app/ 


echo "=================================="
echo "run main application"
echo "=================================="
npm install express
npm install cors
npm install store
npm install body-parser
node server.js
