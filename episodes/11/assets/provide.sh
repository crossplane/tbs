#!/usr/bin/env bash

export BASE64ENCODED_AWS_PROVIDER_CREDS=$(echo -e "[default]\naws_access_key_id = $(aws configure get aws_access_key_id --profile default)\naws_secret_access_key = $(aws configure get aws_secret_access_key --profile default)" | base64  | tr -d "\n")
sed "s/BASE64ENCODED_AWS_PROVIDER_CREDS/$BASE64ENCODED_AWS_PROVIDER_CREDS/g" aws-provider.yaml | kubectl create -f -