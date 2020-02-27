#!/usr/bin/env bash

export AWS_REGION=us-west-2
export AWS_ACCESS_KEY_ID=$(aws configure get aws_access_key_id --profile default)
export AWS_SECRET_ACCESS_KEY=$(aws configure get aws_secret_access_key --profile default)
export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm alpha bootstrap encode-aws-credentials)

clusterawsadm alpha bootstrap create-stack

kustomize build github.com/kubernetes-sigs/cluster-api-provider-aws//config?ref=master | envsubst | kubectl apply -f -
