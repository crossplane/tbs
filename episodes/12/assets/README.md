# Guide

1. Install Crossplane with provider-gcp

```
kubectl create namespace crossplane-system
helm repo add crossplane-master https://charts.crossplane.io/master/
helm search repo crossplane --devel
helm install crossplane --namespace crossplane-system crossplane-master/crossplane --version $version --set clusterStacks.gcp.deploy=true --set clusterStacks.gcp.version=master --devel
```

2. Create GCP `Provider`

Documentation
[here](https://crossplane.io/docs/v0.8/cloud-providers/gcp/gcp-provider.html).

3. Create `GKEClusterClass`

```
kubectl apply -f gkeclusterclass.yaml
```

4. Create `KubernetesCluster` claims

```
kubectl apply -f k8scluster1.yaml
kubectl apply -f k8scluster2.yaml
```

5. Install Linkerd CLI

```
curl -sL https://run.linkerd.io/install | sh
export PATH=$PATH:$HOME/.linkerd2/bin
```

6. Install Linkerd into Target Clusters

*Note: the `kubectl crossplane pack` command is still pending inclusion as part
of the [Crossplane CLI](https://github.com/crossplane/crossplane-cli). Until
[the PR](https://github.com/crossplane/crossplane-cli/pull/47) is merged, it can
be installed by running `make build` thn `make install` on [this
branch](https://github.com/hasheddan/crossplane-cli/tree/pack).*

```
linkerd install | kubectl crossplane pack - --name linkerd1-install --namespace linkerd1 | kubectl apply -f -
linkerd install | kubectl crossplane pack - --name linkerd2-install --namespace linkerd2 | kubectl apply -f -
```

7. Connect to remote Kubernetes clusters and check Linkerd status

```
kubectl get -n linkerd1 secret k8scluster --template={{.data.kubeconfig}} | base64 --decode > remote1.kubeconfig
kubectl get -n linkerd2 secret k8scluster --template={{.data.kubeconfig}} | base64 --decode > remote2.kubeconfig

linkerd check --kubeconfig=remote1.kubeconfig
linkerd check --kubeconfig=remote2.kubeconfig
```

```
linkerd --kubeconfig=remote1.kubeconfig dashboard
linkerd --kubeconfig=remote2.kubeconfig dashboard
```

8. Deploy Workload into Target Cluster

```
curl -sL https://run.linkerd.io/emojivoto.yml | linkerd inject - --kubeconfig=remote1.kubeconfig | kubectl crossplane pack - --name emoji --namespace linkerd1 | kubectl apply -f -

curl -sL https://run.linkerd.io/emojivoto.yml | linkerd inject - --kubeconfig=remote2.kubeconfig | kubectl crossplane pack - --name emoji --namespace linkerd2 | kubectl apply -f -
```

9. View Workload running

```
kubectl get kubernetesapplicationresources -n linkerd1 | grep emoji
kubectl get kubernetesapplicationresources -n linkerd2 | grep emoji
```

If you `kubectl describe` one of the services, there should be a publicly
available IP address where you can access the emoji web app!
