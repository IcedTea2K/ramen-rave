#!/bin/bash

buildOnly=false

while getopts ":h:b" option; do
   case $option in
      h)
         echo "-h: display this help menu"
         echo "-b: build the wasm binaries only"
         exit;;
      b)
         buildOnly=true
         exit;;
     \?)
         exit;;
   esac
done

GOOS=js GOARCH=wasm go build -o ./firefox-extension/bin/main.wasm

# Run the web-extension
if [ $buildOnly != true ]; then
   web-ext run -s ./firefox-extension/
fi
