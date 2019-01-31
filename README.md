 # Kubernetes Deployment Daemon #
 
## Feature ##
1. Create Deploy Daemon will trigger Deployment deploy
2. Make change on Deploy Daemon will impact Deployment change 
3. Support multiple Deployment version online at same time ( share the same virtual service )
4. Support make special version of deployment instance offline 
5. Support trigger deployment with time schedule.
6. Support control percentage of pod from same deployment online or offline  ( TBD )
7. Support control special pod offline from target deployment ( TBD ) 

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
