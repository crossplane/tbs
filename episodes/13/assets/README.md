# Guide

1. Create a `KIND` cluster

```
kind create cluster
```

1. Install `arkade`

```
curl -SLsf https://dl.get-arkade.dev/ | sudo sh
```

1. Install Crossplane using `arkade`

```
arkade install crossplane --helm3
```

1. Install `provider-aws`

```
kubectl apply -f provider-aws.yaml
```

2. Create AWS `Provider`

Documentation
[here](https://crossplane.io/docs/v0.8/cloud-providers/aws/aws-provider.html).

3. Create `S3BucketClass`

```
kubectl apply -f s3bucketclass.yaml
```

4. Install `faas-cli`

```
curl -sL https://cli.openfaas.com | sudo sh
```

5. Install `openfaas` using `arkade`

```
arkade install openfaas
```

6. Port forward dashboard and get credentials

```
kubectl port-forward -n openfaas svc/gateway 8080:8080

PASSWORD=$(kubectl get secret -n openfaas basic-auth -o jsonpath="{.data.basic-auth-password}" | base64 --decode; echo)                                                                                    
echo -n $PASSWORD | faas-cli login --username admin --password-stdin
```

7. Build function

```
<!-- cd func && faas-cli template pull https://github.com/openfaas-incubator/golang-http-template -->

faas-cli build -f tbs.yml --build-arg GO111MODULE=on
```

8. Push image to docker hub

```
faas-cli push -f tbs.yml
```

9. Deploy function

```
faas-cli deploy -f tbs.yml
```

10. Post file

```
curl -X POST -F 'file=@test-image.jpg' http://127.0.0.1:8080/function/upload
```
