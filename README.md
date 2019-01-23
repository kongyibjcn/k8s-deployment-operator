 # Kubernetes Deployment Daemon #

## Generate DeployDaemon Scheme

1. Set Environment Parameter
```
$ ROOT_PACKAGE="github.com/resouer/k8s-controller-custom-resource"
# API Group
$ CUSTOM_RESOURCE_NAME="samplecrd"
# API Version
$ CUSTOM_RESOURCE_VERSION="v1"
```

2. Install k8s.io/code-generator
```
$ go get -u k8s.io/code-generator/...
$ cd $GOPATH/src/k8s.io/code-generator
```
Make user the package path is src/k8s.io/code-generator/

3. Generate Scheme Code
```
$ ./generate-groups.sh all "$ROOT_PACKAGE/pkg/client" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION"
```
Will generate pkg/client pkg/apis and deepcopy code for each object