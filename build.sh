#!/bin/bash

GOOS=js GOARCH=wasm go build -o ./firefox-extension/bin/main.wasm
