#!/bin/bash

rm assets/static/client.js*
cd client
gopherjs build
cd ../assets/static
ln -s ../../client/client.js* .
cd ../..

kick -appPath=$PWD -mainSourceFile=main.go -gopherjsAppPath=client