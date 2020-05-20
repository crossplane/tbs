# Guide

## Create Kind cluster

```
kind create cluster --config kind-config.yaml
```

## Install Crossplane and GCP Provider

Install Crossplane from `alpha` channel
```
kubectl create namespace crossplane-system
helm repo add crossplane-alpha https://charts.crossplane.io/alpha
helm install crossplane --namespace crossplane-system crossplane-alpha/crossplane --set clusterStacks.gcp.deploy=true --set clusterStacks.gcp.version=v0.9.0
```

## Configure GCP credentials

Instructions for setting up your GCP credentials can be found [here](https://crossplane.io/docs/v0.10/cloud-providers/gcp/gcp-provider.html).

## Install Velero

Start by installing the [Velero CLI](https://velero.io/docs/v1.3.2/basic-install/#install-the-cli).

Install Velero into your cluster with AWS provider
```
cat ~/.aws/credentials > cloud-credentials

velero install --provider velero.io/aws --bucket velero-backups --plugins velero/velero-plugin-for-aws:v1.0.1 --bucket crossplane-velero --backup-location-config region=us-east-2 --use-volume-snapshots=false --secret-file=./cloud-credentials
```

## Provision GCP resources using Crossplane

```
kubectl apply -f scenario-1/cloudsqlclass.yaml
kubectl apply -f scenario-1/psql-claim.yaml
```

## Create cluster backup to S3 bucket using Velero

```
velero backup create tbs-16
```

## Destroy Kind cluster

```
kind delete cluster
```

## Create new Kind cluster

```
kind create cluster --config kind-config.yaml
```

## Install Velero

Install Velero into your cluster with AWS provider
```
cat ~/.aws/credentials > cloud-credentials

velero install --provider velero.io/aws --bucket velero-backups --plugins velero/velero-plugin-for-aws:v1.0.1 --bucket crossplane-velero --backup-location-config region=us-east-2 --use-volume-snapshots=false --secret-file=./cloud-credentials
```

## Restore backup in new cluster

```
velero restore create --from-backup tbs-16
```

## A more complex scenario

```
kubectl apply -f scenario-2/install-stack-gcp-sample.yaml
kubectl apply -f scenario-2/install-app-wordpress.yaml
```

```
kubectl apply -f scenario-2/gcp-sample.yaml
kubectl apply -f scenario-2/wordpress.yaml
```


## Create cluster backup to S3 bucket using Velero

```
velero backup create tbs-16-complex
```

## Destroy Kind cluster

```
kind delete cluster
```

## Create new Kind cluster

```
kind create cluster --config kind-config.yaml
```

## Install Velero

Install Velero into your cluster with AWS provider
```
cat ~/.aws/credentials > cloud-credentials

velero install --provider velero.io/aws --bucket velero-backups --plugins velero/velero-plugin-for-aws:v1.0.1 --bucket crossplane-velero --backup-location-config region=us-east-2 --use-volume-snapshots=false --secret-file=./cloud-credentials
```

## Restore backup in new cluster

```
velero restore create --from-backup tbs-16-complex
```