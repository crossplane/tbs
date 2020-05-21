# Guide

This guide is inspired by the "From Zero to K8S to Leafnodes Using Helm" guide
in the [NATS
documentation](https://docs.nats.io/nats-on-kubernetes/from-zero-to-leafnodes).

## Install Crossplane with provider-gcp

Install Crossplane:
```
kubectl create namespace crossplane-system
helm repo add crossplane-alpha https://charts.crossplane.io/alpha
helm install crossplane --namespace crossplane-system crossplane-alpha/crossplane
```

Install `provider-gcp`:
```
kubectl apply -f provider-gcp.yaml
```

Make sure to configure your GCP and account as described in the Crossplane
[installation
guide](https://crossplane.io/docs/v0.11/getting-started/install-configure.html).

## Create New Infrastructure Types

Create `ClusterRole` so Crossplane can manage the new types:
```
kubectl apply -f clusterrole.yaml
```

Create `InfrastructureDefinition` / `InfrastructurePublication` for
`K8sCluster`:
```
kubectl apply -f k8scluster.yaml
```

## Create Compositions for GCP

```
kubectl apply -f gke.yaml
```

## Create K8sClusterRequirements

```
kubectl apply -f req-central.yaml
kubectl apply -f req-west.yaml
```

You can view any managed resources that are provisioned by running:
```
kubectl get managed
```

## Add NATS Helm Chart

```
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm repo update
```

## Install NATS with Leaf Node Configuration


Get `kubeconfig` for each cluster:
```
kubectl get secret k8s-conn-central --template={{.data.kubeconfig}} | base64 --decode > central.kube
kubectl get secret k8s-conn-west --template={{.data.kubeconfig}} | base64 --decode > west.kube
```

Create ngs-creds in each cluster.

> See more in the [NGS Getting Started Guide](https://synadia.com/ngs/signup).

```
kubectl --kubeconfig=central.kube create secret generic ngs-creds --from-file $HOME/.nkeys/creds/synadia/TBS/TBS.creds
kubectl --kubeconfig=west.kube create secret generic ngs-creds --from-file $HOME/.nkeys/creds/synadia/TBS/TBS.creds
```

Deploy NATS leafnodes in each cluster:
```
helm --kubeconfig=central.kube install nats nats/nats -f nats.yaml
helm --kubeconfig=west.kube install nats nats/nats -f nats.yaml
```

Start reply servers on `nats-box` in each cluster:
```
kubectl --kubeconfig=central.kube exec -it nats-box -- nats-rply hello "cluster-central"
kubectl --kubeconfig=west.kube exec -it nats-box -- nats-rply hello "cluster-west
```

Make request from local machine:
```
docker run --rm -it -v ~/.nkeys/creds/synadia/TBS/TBS.creds:/creds/TBS.crds synadia/nats-box:latest
$ nats-req -creds /creds/TBS.creds hello "From Local!"
```

You should get a response back from whatever region you are closer to:
```
cluster-central
```

Now stop the reply server that you got the response from, and send another
request:
```
$ nats-req -creds /creds/TBS.creds hello "From Local again!"
```

You should get a response back from the other region:
```
cluster-west
```
