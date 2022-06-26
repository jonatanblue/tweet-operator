# TweetOperator

This repo is intended as an example project, showing how to write a custom Kubernetes controller, aka operator. I mainly want to show two things: writing a simple operator without all the bells and whistles is actually pretty straight-forward - as long as you get the code generation right - and that a Kubernetes controller can do pretty much anything: [manage DaemonSets](https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/daemon/daemon_controller.go), [order pizza](https://github.com/rudoi/cruster-api) and tweet.

The TweetOperator posts a tweet for each Tweet custom resource created in the cluster, and posts back status information about the tweet: likes, retweets, etc.

## Setup

Go to https://developer.twitter.com, set up a developer account and fill out the form to apply for **Elevated access**.

Access needed for the app:
* `tweet.read`
* `tweet.write`
* `users.read`

Add the credentials to `env.local`:

```
TWITTER_USERNAME=<some username>
CONSUMER_KEY=<redacted>
CONSUMER_SECRET=<redacted>
ACCESS_TOKEN=<redacted>
ACCESS_TOKEN_SECRET=<redacted>
```

## Development

### Run locally

Export environment variables

```
. env.local
```

Run the operator locally:

```
go run main.go
```

### Run in a cluster

Build Dockerimage

```
docker build --tag tweet-operator:v1 .
```

Set up a KIND cluster for local testing.

```
kind create cluster
```

Load image into kind cluster:

```
kind load docker-image tweet-operator:v1
```

Create the Tweet custom resource:

```
kubectl create -f manifests/hello_world_tweet.yaml
```

Export environment variables

```
. env.local
```

Create secret object for the Twitter API credentials:

```
kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: twitter-credentials
  namespace: default
type: Opaque
data:
  TWITTER_USERNAME: $(printf "$TWITTER_USERNAME" | base64)
  CONSUMER_KEY: $(printf "$CONSUMER_KEY" | base64)
  CONSUMER_SECRET: $(printf "$CONSUMER_SECRET" | base64)
  ACCESS_TOKEN: $(printf "${ACCESS_TOKEN}" | base64)
  ACCESS_TOKEN_SECRET: $(printf "$ACCESS_TOKEN_SECRET" | base64)
EOF
```

Create operator Deployment:

```
kubectl apply -f manifests/operator.yaml
```

## Appendix 1: Code generation

This bit is for your reference, for when you write your own operator. I have structured the commits to split up making the blueprint (the first three files in the `pgk/apis` folder) from the code generation.

### Set up your GOPATH

I am using these folders:

* GOPATH: `/Users/jonatan/go`
* This project: `/Users/jonatan/go/src/github.com/jonatanblue/tweet-operator`
* code-generator: `/Users/jonatan/go/src/k8s.io/code-generator` (cloned from https://github.com/kubernetes/code-generator)

### Run code generator

```
$ codegen_path=/Users/jonatan/go/src/k8s.io/code-generator
$ "${codegen_path}"/generate-groups.sh all github.com/jonatanblue/tweet-operator/pkg/client github.com/jonatanblue/tweet-operator/pkg/apis example.com:v1 --go-header-file "${codegen_path}"/hack/boilerplate.go.txt
Generating deepcopy funcs
Generating clientset for example.com:v1 at github.com/jonatanblue/tweet-operator/pkg/client/clientset
Generating listers for example.com:v1 at github.com/jonatanblue/tweet-operator/pkg/client/listers
Generating informers for example.com:v1 at github.com/jonatanblue/tweet-operator/pkg/client/informers
```

### Generate YAML for registering CRD

First clone the repo for controller-tools and build the binary:

```
git clone git@github.com:kubernetes-sigs/controller-tools.git
cd controller-tools/cmd/controller-gen
go build
```

Then, from the root of your project, run the binary:

```
$ ${path_to_controller_gen}/controller-gen paths=github.com/jonatanblue/tweet-operator/pkg/apis/example.com/v1 crd:crdVersions=v1 output:crd:artifacts:config=manifests
```

This will generate a yaml file in `manifests/`. Use it to register the CRD in the cluster:

```
kubectl create -f manifests/example.com_tweets.yaml
```
