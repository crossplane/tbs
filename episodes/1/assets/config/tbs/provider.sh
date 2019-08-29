#!/usr/bin/env bash

export BASE64ENCODED_TOKEN=$(base64 token.txt | tr -d "\n")
sed "s/BASE64ENCODED_TOKEN/$BASE64ENCODED_TOKEN/g" provider.yaml | kubectl create -f -