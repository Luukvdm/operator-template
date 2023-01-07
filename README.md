# Operator Template
A repository that can be easily used as a base for a Kubernetes operator.
I created this repository because I felt like the files that projects like [Operator-SDK](https://sdk.operatorframework.io/) 
and [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) generate are to bloated and hard to understand.

## Getting Started

### Building and installing

Build the project and install the helm chart for the current context:
```sh
make install
```

Install on a kind cluster, beside doing a normal install, this also loads the image onto the cluster:
```sh
make kind-install
```

### Code Generation

The custom resource definition yaml, role yaml and Go copy functions for the resource are generated.
You can use these make goals to (re)generate them:
```sh
make install-tools
make generate
```

## Dependencies

### Manager
The [controller-runtime manager](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/manager#Manager) is used to 
easily create the controller for the custom resource. This replaces a lot of boilerplate code that allows us to focus on
the logic for our custom resource.

### Generation
This project uses the [controller-tools controller-gen](https://github.com/kubernetes-sigs/controller-tools) from Kubernetes to generate some files.
The files that are generated are the custom resource yaml, role yaml, and the Go copy implementations for the resource.   
By generating the Yaml files the Go code becomes the single source of truth about the custom resource definition properties.

### CLI Interface
Like most Go projects, we use [spf13/cobra](https://github.com/spf13/cobra) to create the CLI interface.
Right now the CLI code is relatively basic, but this could expand the more your project grows.  
And other then that, using Cobra just makes your life easier.

### Logging
The manager from controller runtime makes use of [logr](https://github.com/go-logr/logr) abstraction for loggers.
There are many implementations for this abstraction like [glogr (golog)](https://github.com/go-logr/glogr), 
[klogr (logger from the Kubernetes project)](https://git.k8s.io/klog/klogr),
[stdr (the standard logger library)](https://github.com/go-logr/stdr) and many more.  
This template uses the [zapr implementation](https://github.com/go-logr/zapr) because of it's performance and my preferences.

