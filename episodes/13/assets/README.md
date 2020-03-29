# Guide

1. Create a `KIND` cluster

```
kind create cluster
```

2. Install `arkade`

```
curl -SLsf https://dl.get-arkade.dev/ | sudo sh
```

3. Install Crossplane using `arkade`

```
arkade install crossplane --helm3
```

4. Install `provider-aws`

```
kubectl apply -f provider-aws.yaml
```

5. Create AWS `Provider`

Documentation
[here](https://crossplane.io/docs/v0.8/cloud-providers/aws/aws-provider.html).

6. Create `S3BucketClass`

```
kubectl apply -f s3bucketclass.yaml
```

7. Install `faas-cli`

```
curl -sL https://cli.openfaas.com | sudo sh
```

8. Install `openfaas` using `arkade`

```
arkade install openfaas
```

9. Port forward dashboard and get credentials

```
kubectl port-forward -n openfaas svc/gateway 8080:8080 &

PASSWORD=$(kubectl get secret -n openfaas basic-auth -o jsonpath="{.data.basic-auth-password}" | base64 --decode; echo)                                                                                    
echo -n $PASSWORD | faas-cli login --username admin --password-stdin

echo -n $PASSWORD | xclip -sel clip
```

10. Create `Bucket`

```
kubectl apply -f bucket.yaml
```

11. Build function

```
cd func

<!-- faas-cli template pull https://github.com/openfaas-incubator/golang-http-template -->

faas-cli build -f tbs.yml --build-arg GO111MODULE=on --parallel 3
```

12. Push image to docker hub

```
faas-cli push -f tbs.yml
```

13. Deploy function

```
faas-cli deploy -f tbs.yml
```

14. Post file

```
cd ..

curl -X POST -F 'file=@test-image.jpg' http://127.0.0.1:8080/function/upload
```

15. View UI

Go to `http://127.0.0.1:8080/function/upload`
