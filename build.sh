#!/bin/bash

rm assets/static/client.js*
cd client
gopherjs build
cd ../assets/static
cp ../../client/client.js* .
cd ../..

packr2 build
packr2 clean