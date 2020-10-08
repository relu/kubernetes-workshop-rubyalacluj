# Kubernetes Introduction Workshop

## Prerequisites

Install the following tools on your machine:

### Kind

Please follow the guide
[here](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) to install
kind.  To set up the local cluster run the following command:

```
$ kind create cluster --name workshop --config manifests/00-kind-config.yaml
```

The configuration contains some settings that allow us to set up ingresses
later on.

### Kubectl

Follow [the guide](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

### Helm

Install using the guide [here](https://github.com/helm/helm#install)

### stern

Install from [here](https://github.com/wercker/stern)

## Clone this repository

In order for us to work through the examples more easily make sure to clone
this repository and switch your current working path to its location.

All commands involving files should run from the repository base directory.

## Create a namespace

A namespace represents a means of partitioning a kubernetes cluster. Most
Kubernetes objects are namespaced, with the exception of namespaces themselves,
persistent volumes and cluster nodes.

We'll use this namespace to create resources inside it and once we're done
we'll be able to quickly clean everything up by deleting the namespace which
will also delete all resources contained in it.

```
$ kubectl create namespace kubealacluj
```

Now that we have your namespace, let's see if we can list it out:

```
$ kubectl get ns
```

_NOTE: Some resources have short names, in the case of namespaces, the short
name is `ns`._

You should now be able to identify it in the list and it should have a status
of active.

Now let's inspect this namespace resource and see how it looks like:

```
$ kubectl describe ns kubealacluj
```

This should show us a more detailed view of our namespace's current state.

_NOTE: All k8s resources can be described using the `kubectl describe <resource
type>` command._

As mentioned in the presentation earlier, all k8s objects are created through
manifest files. Namespaces are no exception, even though they can be created
ad-hoc via the `kubectl create` command that doesn't mean they are not backed
by a declarative manifest.

Let's see how the namespace manifest looks like:

```
$ kubectl get -o yaml ns kubealacluj
```

_NOTE: As with other commands described earlier, you can use `kubectl get` on
all k8s objects._

What if we want to edit the resource in place? We can surely do that!

```
$ kubectl edit ns kubealacluj
```

This will bring up your editor and give you the posibility to edit the
resource directly. Once saved, it will be applied in kubernetes automatically.

Ok, so now that we have our workspace ready, let's set our current kubectl
context to make use of it so that all future interactions will use our
namespace.

```
$ kubectl config set-context --current --namespace kubealacluj
```

_NOTE: By default, contexts will have the namespace set to the `default`
namespace. The `default` namespace is created automatically by k8s and you
should have already seen it in the list of namespaces._

## Create a pod

Alright, let's get to some more serious matters. In the `apps/` directory you
should be able to spot three small example applications built in three
different programming languages: python, go and of course, ruby.
We'll use these app images:

- aurelcanciu/example-app-python
- aurelcanciu/example-app-ruby
- aurelcanciu/example-app-go

_NOTE: Feel free to inspect these, change them and build your own images._

Now let's create our pod:

```
$ kubectl apply -f manifests/01-pod.yaml
```

And now let's see if it's running:

```
$ kubectl get po
```

_NOTE: The `metadata.name` field in the manifest is the name that our resource
will have assigned, this name should always be unique in the scope of the
namespace._


### We have a pod, what next? How can we interact with it?

Luckly we can actually port forward a localhost port the pod, here's how:

```
$ kubectl port-forward pod/web-app 3000
```

This basically binds localhost:3000 to our pod's 3000 port. You can now access
http://localhost:3000

_NOTE: Only TCP is currently supported for port forwarding and it can work for
Pod, ReplicaSet, Deployment and Service resources._

### What else can I do with this pod?

You can actually exec into the pod (if there's a shell available in the image
of course). Our pod container's image is based on
[alpine](https://alpinelinux.org/) so we should be able to use `ash`, its
default shell.

```
$ kubectl exec -ti web-app ash
```

This should bring up the shell prompt and voila, you're in!

## Create a service

Now that we have a pod, let's create a service object that would act as a proxy
for our pod.

```
$ kubectl apply -f manifests/02-service.yaml
```

And we should see our newly created service:

```
$ kubectl get svc

```

_NOTE: Our service is of type NodePort, this service type will create a unique
host port mapping which can be accessed from outside the cluster if needed.
Other types are ClusterIP (no host port mapping, only reachable from inside the
cluster network) and LoadBalancer (provisions a cloud provider load balancer
for our service)._

Let's port forward to our service:
```
$ kubectl port-forward svc/web-app 3000
```

Now http://localhost:3000 will forward traffic to our service which will proxy
it to our pod.

The way this works is through using label selectors. As you can see, our pod
has an `app` label while our service resource uses the same label key-value
pair as a `selector`. This means that the service will self-discover all pods
labeld with the same values as its selector key-value pairs.

## Replication

So what if we want to scale-out our application for resilience and high
availability? Sounds like we need to create more of them pods... Here's where
the ReplicaSet comes in handy.

```
$ kubectl apply -f manifests/03-replicaset.yaml
```

This will create our replicaset configured to replicate 3 pods. Let's check it
out:

```
$ kubectl get rs
$ kubectl get pod
```

As you can see we have three pods running right now. But wait, why are they
three when we instructed the RestfulSet to bring up three and we had one
already created statically from before? It's simple, the first pod we created
has the same label as the pods managed by the ReplicaSet, this means that the
ReplicaSet will count that in as well as long as its selector will match it.

The ReplicaSet will always maintain a fixed number of pods as per its
configuration. Think of it as a means to group multiple pods of the same type
together. If one pod is terminated, a new pod will be crated in its place.

Let's delete our initially created static pod and see what happens.

```
$ kubectl delete -f manifests/01-pod.yaml
```

_NOTE: We used the manifest file to delete the pod resource, however supplying
the pod name would work as well, you don't need the manifest to delete
resources._

Let's check the pods now and see what's going on:

```
$ kubectl get po
```

Looks like we have 3 pods now and all of them are controlled by the ReplicaSet.
If we describe the existing service, we'll see that it will have 3 endpoints
corresponding to the 3 pods:

```
$ kubectl describe svc web-app
```

## Deployments

How do we do a rolling update of our application right now? Well, it's not
easy. You would need to create a new replicaset for the new version and then
scale in the old replicaset and clean it up yourself, which is a bit annoying.
That's why Kubernetes introduced the Deployment resource.

```
$ kubectl apply -f manifests/04-deployment.yaml
```

_NOTE: If you look at its manifest, the deployment looks identical to the
ReplicaSet, aside from the `kind` atttribute that is. The reason is that
Deployments are controllers for replicasets, they basically work by creating,
scaling and removing replicasets._

Our deployment should be now created:
```
$ kubectl get deployments
```

If we look at the pods again you'll notice nothing really changed, same reason
as for why the replicaset didn't create an extra before and that's because the
selector matches the existing pods which prevent creating new ones.


### Rolling update

Now let's make a change to our deployment so that an update takes place:

```
$ kubectl set image deployment/web-app app=aurelcanciu/example-app-python --record
```

Now, if we look at the pods, we'll see that the rolling update process started
and new pods are being launched. They are being replaced one by one by the
deployment logic.

Let's check out the replicasets:

```
$ kubectl get rs
```

You should see two replicasets, one is the old one which has desired size set
to 0 and the other is the new one which has desired size set to 3.

### Scaling a deployment

If we need to manually scale our deployment, this can be done simply and
quickly:

```
$ kubectl scale deployment web-app --replicas 5
```

Let's update the deployment once again with the go image:

```
$ kubectl set image deployment/web-app app=aurelcanciu/example-app-go --record
```

_NOTE: We're updating the image but we can actually update just the image tag
and that's basically how things would be done in most situations._

### Useful deployment related commands

We can check the rollout status of a deployment by running:

```
$ kubectl rollout status deployment web-app
```

To see the deployment rollout history you can run:

```
$ kubectl rollout history deployment web-app
```

This will give us a list of revisions with their respective change cause.

### Rollback

To rollback to the previous deployment revision you can simply run:

```
$ kubectl rollout undo deployment web-app
```

If you need to rollback to a specific revision then use the `--to-revision`
argument:

```
$ kubectl rollout undo deployment web-app --to-revision=1
```

## Ingress

Now that we have a deployment running for our app we're ready to expose it to
the external world. In order to do this we'll need to create an Ingress object.

Normally an ingress object maps 1:1 with an external Load Balancer serivce,
since we're running locally and not in the cloud we'll need to install an
ingress controller. To install the nginx ingress controller component run:

```
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
```

If you look in the `ingress-nginx` namespace, you'll see there's a new pod
running there with a name prefix `ingress-nginx-controller-`:

```
$ kubectl get pod -n ingress-nginx
```

Now we're ready to create our Ingress:

```
$ kubectl apply -f manifests/05-ingress.yaml
```

You'll need to wait some time until the ingress is ready. Can check it out by
running:

```
$ kubectl describe ingress web-app
```

We should now be able to access the new ingress load balancer on localhost.

## Resource management

In order to get metrics working we need to install the metrics-server, it's
really simple if we use helm to do this (more on this a bit later).

First we need to add the `stable` repository helm will use to grab the
metrics-server chart:

```
$ helm repo add stable https://kubernetes-charts.storage.googleapis.com
```

Then we can install it:


```
$ helm install metrics-server stable/metrics-server -n kube-system --set 'args[0]=--kubelet-insecure-tls'
```

Alternatively, we can install metrics-server via kubectl using the manifest:

```
$ kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.3.7/components.yaml
```

And apply some changes to the deployment so it works with our kind setup:

```
$ kubectl edit deployments.apps -n kube-system metrics-server
```

Add to container args:

```
args:
  - --kubelet-insecure-tls
  - --kubelet-preferred-address-types=InternalIP
```

To get a sense of how much resources a pod consumes we can use the following
command:

```
$ kubectl top pod
```

Now that we know how much resources our pods are using, we can tweak the pod
configuration by adding a resources section:

```yaml
resources:
  requests:
    cpu: 0.1
    memory: 50Mi
  limits:
    cpu: 0.25
    memory: 100Mi
```

Let's try to understand these a bit better:

- `requests` refers to the resources the cluster needs to have available on a
  particular node in order to schedule the pod, we can refer to it as reservation
  as well.
- `limits` refers to the maximum amount of CPU and memory a particular pod
  container may consume and it represents a hardly enforced limit, meaning that
  the pod will be throttled on CPU utilizationand can be OOM killed if it
  exeeds the memory limit.

### CPU

CPU is specified in units of cores: 1 CPU core = 1 cpu unit = 1000m (milli)cpu units

### Memory

Memory is specified in units of bytes: 1Mi = 1 mebibyte =  1024 * 1024 bytes

### Apply the resources configuration

```
$ kubectl apply -f manifests/06-deployment-with-resources.yaml
```

_NOTE: This configuration is at the level of a single pod container. Since you
can run multiple containers in a pod, each can have its own specific resource
requirements. In a lot of cases, sidecar containers will require a lot less
resources than the application they side with._

## Configuration

Let's apply a new manifest file:

```
$ kubectl apply -f manifests/07-configuration.yaml
```

_NOTE: As mentioned earlier, a kubernetes manifest file can contain multiple
object definitions with a file delimiter used between them `---`._

### ConfigMaps

Let's see how ConfigMaps can be used to manage application configuration:

The manifest has three configmaps defined, we can see they got created by
listing them:

```
$ kubectl get cm
```

Now let's see how the ConfigMaps are used by our pods. We can describe one of
the pods to see:

```
$ kubectl describe pod -l app=web-app
```

_NOTE: Some kubectl subcomands such as `get`, `describe`, `delete` etc. accept
a label selector `--selector` or `-l` argument which allows us to filter out
the objects we want to retrieve._

We can now see our pods use the three configmaps in three different ways:
1. Environment
2. Environment Variables from
3. Mounts (volumes)

If we exec into one of the pod containers we'll be able see the environment
variables and check the `/tmp/config` directory which should have a `file.txt`
in it with the contets from our ConfigMap.

### Secrets

Inside the manifests secret values are stored as base64 encoded values. To
encode values you can use the following command:

```
$ echo -n 'secret value' | base64
```

_NOTE: The `type: Opaque` means that the secret contains arbitrary
non-structured data. Kubernetes has other secret types such as ServiceAccount
tokens or secrets used to authenticate private image repositories which are
constrained to a particular schema._

Secrets can be used the same way ConfigMaps are within a pod and the values
will be exposed to the pod containers decrypted.

## Logging

We are able to see the logs our pod applications produce by using the `kubectl
logs` subcommand.

```
$ kubectl logs -f -l app=web-app
```

For better results, you can use `stern` :)

## Kubernetes dashboard

Install the kubernetes dashboard via manifest:

```
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.4/aio/deploy/recommended.yaml
```

And create an admin user:

```
$ kubectl apply -f manifests/08-dashboard-admin.yaml
```

Before we can use the dashboard we need a token to authenticate:

```
$ kubectl -n kubernetes-dashboard describe secret $(kubectl -n kubernetes-dashboard get secret | grep admin-user | awk '{print $1}') | awk '$1=="token:"{print $2}'
```

And now run the kubectl proxy:

```
$ kubectl proxy
```

Now you can access the dashboard at:
http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/

## Helm

Helm is kind of like a package manager for Kubernetes. Simply put, it enables
users to define manifest files and create a unified package which we call a
Chart.

Helm is also able to retrieve charts from dedicated repositories and it also has
dependency management which means that a chart can rely on multiple subcharts.

A large collection of community maintained charts can be found here:
https://github.com/helm/charts

### Our own chart

In the `helm/` directory you can find an example chart created based on the
manifests we've worked through until now. Our helm chart also has a dependency
on `redis`.

Let's instll our first helm release of the `example-app` chart into the
cluster:

```
$ helm install example-app-ruby ./helm/example-app --set image.repository=aurelcanciu/example-app-ruby --wait
```

Now we can see that our release was created by listing the available releases:

```
$ helm ls
```

_NOTE: In Helm v3, releases are namespaced._

What if we want to upgrade our release? Let's create the ingress which was
previously disabled:

```
$ helm upgrade example-app-ruby ./helm/example-app --reuse-values --set ingress.enabled=true
```

Now we should see that our release revision was incremented and we should have
a new ingress object.

Let's create releases for our other apps:

```
$ helm install example-app-python ./helm/example-app --set image.repository=aurelcanciu/example-app-python -f ./helm/overrides.yaml --wait
$ helm install example-app-go ./helm/example-app --set image.repository=aurelcanciu/example-app-go -f ./helm/overrides.yaml --wait
```

As you can see, we've used the same chart to create two new releases with a
different images. We also used a file to override some chart values instead of
providing them via the command line. I imagine you can by now figure out how
useful can helm be when it comes to application deployment management.

### Uninstalling a release

This is as simple as running:

```
$ helm uninstall example-app-ruby
```

That's it!

## Cleanup

It's done, so let's clean up the mess :)

```
$ kind delete cluster --name workshop
```

Au revoir!
