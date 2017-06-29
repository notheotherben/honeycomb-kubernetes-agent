![](static/hero.png)

## Getting started

If you don't already have a Kubernetes cluster handy, there are a couple of ways to get one. You can use [minikube](https://kubernetes.io/docs/getting-started-guides/minikube/) to run a local one-node cluster in a VM, or Heptio's [AWS quickstart](https://aws.amazon.com/quickstart/architecture/heptio-kubernetes/) to bring up a cluster in AWS.

You should also install `kubectl`: https://kubernetes.io/docs/tasks/tools/install-kubectl/

Next, clone this repository:
```
git checkout https://github.com/honeycombio/honeycomb-kubernetes-agent
cd honeycomb-kubernetes-agent
git checkout emfree.tutorial
cd docs
```

## Deploying a sample application on Kubernetes

We're going to deploy, scale, and canary a minimal "hello world" application.

First, get yourself your own namespace:

```
export $MY_NAMESPACE=<YOUR_NAME_HERE>
kubectl create namespace $MY_NAMESPACE
kubectl config set-context $(kubectl config current-context) --namespace=$MY_NAMESPACE
```

Now create a [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) of the "hellow world v1" application:

```
kubectl apply -f tutorial/helloworld-1.yaml
```

You can see that you've created some resources:
```
kubectl get pods
kubectl get deployments
```

But we have no way of talking to them. Let's remedy that
```
kubectl expose deployment helloworld-v1 --type=ClusterIP --name=helloworld-v1-service
kubectl get services
```

Now we can talk to our "Hello World" service in a well-defined way from inside
the cluster:
```
kubectl run -it curl --image=tutum/curl
# Now you have a shell in a container in the cluster
$ curl helloworld-v1-service:5000/hello
```

But not from outside the cluster. There are a couple of different ways to address that. Let's expose our service to the outside world:

```
kubectl create -f tutorial/helloworld-svc.yaml
kubectl config set-context $(kubectl config current-context) --namespace=default
kubectl apply -f tutorial/istio/istio-rbac-beta.yaml
kubectl apply -f tutorial/istio/istio.yaml
```

Aaand let's also capture some events:
```
cd ..
kubectl create configmap honeycomb-agent-config --from-file=config.yaml --namespace=kube-system
kubectl create -f ./honeycomb-agent-ds.yml
```

