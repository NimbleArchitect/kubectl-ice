# kubectl-ice
A kubectl plugin that lets you can see the running configuration of all containers
 that are running inside pods, I created it so I could peer inside the pods and see
 the details of containers (sidecars) running in a pod and then extended it so all
 containers could be viewed at once.

ice lists useful information about the (sidecar) containers present inside a
 pod, useful for trouble-shooting multi container issues. You can view volume, 
 image, port and executable configurations, along with current cpu and memory
 metrics all at the container level (requires metrics server)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nimblearchitect/kubectl-ice)
![GitHub](https://img.shields.io/github/license/NimbleArchitect/kubectl-ice)
![Github All Releases](https://img.shields.io/github/downloads/NimbleArchitect/kubectl-ice/total.svg?color=blue)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/NimbleArchitect/kubectl-ice/release)

## Features:
* Only uses read permissions, no writes are called
* Lists the containers in all pods in the current namespace and context
* Selectors work just like they do with the standard kubectl command
* Sortable output columns 
* Can list all containers from all pods across all namespaces
* Exclude rows from output using the match flag, useful to exclude containers with low memory or cpu usage
* List only cpu and memory results that dont fall within range using the oddities flag

# Installation

## Install using krew

```
$ kubectl krew install ice
```
update with 
```
$ kubectl krew upgrade ice
```
dont have krew? check it out here [https://github.com/GoogleContainerTools/krew](https://github.com/GoogleContainerTools/krew)

## Install from binary
- download the required binary from the release page
- unzip and copy the kubectl-ice file to your path
- run ```kubectl-ice help``` to check its working

## Install from Source

clone and build the source using the following commands
```shell
git clone https://github.com/NimbleArchitect/kubectl-ice.git
cd kubectl-ice
make bin
```
then copy ./bin/kubectl-ice to somewhere in your path and run ```kubectl-ice help``` to check its working

## Usage
if kubectl-ice is in your path you can replace the command ```kubectl-ice``` with ```kubectl ice``` (remove the dash) to
 make it feel more like a native kubectl command, this also works if you have kubectl set as an alias, for example
 if k is aliased to kubectl you can type ```k ice subcommand``` instead of ```kubectl-ice subcommand```


The following commands are available for `kubectl-ice`
```
kubectl-ice command    # retrieves the command line and any arguments specified at the container level
kubectl-ice cpu        # return cpu requests size, limits and usage of each container
kubectl-ice help       # Shows the help screen
kubectl-ice image      # list the image name and pull status for each container
kubectl-ice ip         # list ip addresses of all pods in the namespace listed
kubectl-ice memory     # return memory requests size, limits and usage of each container
kubectl-ice ports      # shows ports exposed by the containers in a pod
kubectl-ice probes     # shows details of configured startup, readiness and liveness probes of each container
kubectl-ice restarts   # show restart counts for each container in a named pod
kubectl-ice status     # list status of each container in a pod
kubectl-ice volumes    # list all container volumes with mount points
```

ice also supports all the standard kubectl flags in addition to:
```
Flags:
  -A, --all-namespaces                 List containers from pods in all namespaces
  -c, --container string               Container name. If set shows only the named containers
      --context string                 The name of the kubeconfig context to use
      --match string                   Filters out results, comma seperated list of COLUMN OP VALUE, where OP can be one of ==,<,>,<=,>= and != 
  -n, --namespace string               If present, the namespace scope for this CLI request
  -l, --selector string                Selector (label query) to filter on
```
select subcommands also support the following flags
```
Flags:
  -p, --previous         show previous state
  -r, --raw              show raw uncooked values
      --sort string      Sort by column
      --oddities         show only the outlier rows that dont fall within the computed range (requires min 5 rows in output)
```
all flags are optional, see usage instructions and examples for more info

## Examples
Some examples are listed below but full [usage instructions](https://github.com/NimbleArchitect/kubectl-pod/blob/main/docs/usage.md) and [examples](https://github.com/NimbleArchitect/kubectl-pod/blob/main/docs/examples.md) can be found in the [docs directory](https://github.com/NimbleArchitect/kubectl-pod/blob/main/docs)

### Single pod info
Shows the currently used memory along with the configured memory requests and limits of all containers (side cars) in the pod named web-pod
``` shell
$ kubectl-ice memory web-pod
CONTAINER    USED  REQUEST  LIMIT  %REQ    %LIMIT
app-watcher  0.29M 1M       512M   29.08   0.06
app-broken   0     1M       512M   -       -
myapp        7.61M 1M       256M   760.63  2.97

```
### Using labels
using labels you can search all pods that are part of a deployment where the label app matches demoprobe and list selected information about the containers in each pod, this example shows the currently configured probe information and gives details of configured startup, readiness and liveness probes of each container
``` shell
$ kubectl-ice probes -l app=demoprobe
PODNAME                      CONTAINER     PROBE     DELAY  PERIOD  TIMEOUT  SUCCESS  FAILURE  CHECK    ACTION
demo-probe-76b66d5766-j9wnm  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-j9wnm  web-frontend  liveness  5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-j9wnm  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/
demo-probe-76b66d5766-ksn5t  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-ksn5t  web-frontend  liveness  5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-ksn5t  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/

```
### Container status
most commands work the same way including the status command which also lets you see which container(s) are causing the restarts and by using the optional --previous flag you can view the containers previous exit code
``` shell
$ kubectl-ice status -l app=myapp --previous
T  PODNAME  CONTAINER    STATE       REASON  EXIT-CODE  SIGNAL  TIMESTAMP                      MESSAGE
S  web-pod  app-broken   Terminated  Error   1          0       2022-05-23 10:59:49 +0100 BST  -
S  web-pod  app-watcher  Terminated  Error   2          0       2022-05-23 10:56:58 +0100 BST  -
S  web-pod  myapp        Terminated  Error   137        0       2022-05-21 18:51:29 +0100 BST  -
I  web-pod  app-init     -           -       -          -       -                              -

```
### Advanced labels
return memory requests size and limits of each container where the pods have an app label that matches useoddcpu and the container name is equal to nginx
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-459t8  web-frontend  2m    1m       1000m  155.10    0.16
demo-odd-cpu-5f947f9db4-6mlk9  web-frontend  2m    1m       1000m  149.30    0.15
demo-odd-cpu-5f947f9db4-7xcqw  web-frontend  2m    1m       1000m  134.58    0.13
demo-odd-cpu-5f947f9db4-8fc4c  web-frontend  2m    1m       1000m  145.22    0.15
demo-odd-cpu-5f947f9db4-9x5mb  web-frontend  2m    1m       1000m  142.69    0.14
demo-odd-cpu-5f947f9db4-bxchg  web-frontend  96m   1m       1000m  9567.21   9.57
demo-odd-cpu-5f947f9db4-fsccd  web-frontend  2m    1m       1000m  146.66    0.15
demo-odd-cpu-5f947f9db4-gtlcl  web-frontend  2m    1m       1000m  139.99    0.14
demo-odd-cpu-5f947f9db4-j882g  web-frontend  2m    1m       1000m  152.87    0.15
demo-odd-cpu-5f947f9db4-mqwnd  web-frontend  2m    1m       1000m  137.31    0.14
demo-odd-cpu-5f947f9db4-qh7gk  web-frontend  2m    1m       1000m  180.17    0.18
demo-odd-cpu-5f947f9db4-rcxjq  web-frontend  2m    1m       1000m  149.45    0.15
demo-odd-cpu-5f947f9db4-rrj7c  web-frontend  2m    1m       1000m  154.26    0.15
demo-odd-cpu-5f947f9db4-rtxlm  web-frontend  105m  1m       1000m  10461.33  10.46
demo-odd-cpu-5f947f9db4-xs2gs  web-frontend  2m    1m       1000m  155.57    0.16
demo-odd-cpu-5f947f9db4-zx5c8  web-frontend  2m    1m       1000m  140.37    0.14

```
### Odditites and sorting
given the listed output above the optional --oddities flag picks out the containers that have a high cpu usage when compared to the other containers listed we also sort the list in descending order by the %REQ column
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend --oddities --sort '!%REQ'
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-rtxlm  web-frontend  105m  1m       1000m  10461.33  10.46
demo-odd-cpu-5f947f9db4-bxchg  web-frontend  96m   1m       1000m  9567.21   9.57

```
### Pod volumes
list all container volumes with mount points
``` shell
$ kubectl-ice volumes web-pod
CONTAINER    VOLUME                 TYPE       BACKING           SIZE  RO    MOUNT-POINT
app-init     kube-api-access-4h6h2  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-watcher  kube-api-access-4h6h2  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-broken   kube-api-access-4h6h2  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
myapp        app                    ConfigMap  app.py            -     false /myapp/
myapp        kube-api-access-4h6h2  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount

```


## Contributing

All feedback and contributions are welcome, if you want to raise an issue or help with fixes or features please [raise an issue to discuss](https://github.com/NimbleArchitect/kubectl-ice/issues)


## License
Licensed under Apache 2.0 see [LICENSE](https://github.com/NimbleArchitect/kubectl-pod/blob/main/LICENSE)
