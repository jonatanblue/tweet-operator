# Setup

Add the following to `env.local`:

```
export TWITTER_API_KEY="redacted"
export TWITTER_API_SECRET="redacted"
export TWITTER_BEARER_TOKEN="redacted"
```

# Generate code

## Set up your GOPATH

I am using these folders:

* GOPATH: `/Users/jonatan/go`
* This project: `/Users/jonatan/go/src/github.com/jonatanblue/tweet-operator`
* code-generator: `/Users/jonatan/go/src/k8s.io/code-generator`

## Run code generator

```
$ codegen_path=/Users/jonatan/go/src/k8s.io/code-generator
$ "${codegen_path}"/generate-groups.sh all github.com/jonatanblue/tweet-operator/pkg/client github.com/jonatanblue/tweet-operator/pkg/apis example.com:v1 --go-header-file "${codegen_path}"/hack/boilerplate.go.txt
Generating deepcopy funcs
Generating clientset for example.com:v1 at github.com/jonatanblue/tweet-operator/pkg/client/clientset
Generating listers for example.com:v1 at github.com/jonatanblue/tweet-operator/pkg/client/listers
Generating informers for example.com:v1 at github.com/jonatanblue/tweet-operator/pkg/client/informers
```
