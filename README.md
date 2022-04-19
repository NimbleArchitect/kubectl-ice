# kubectl-ice
A kubectl plugin that lets you can see the running configuration of all containers
 that are running inside pods, I created it so I could peer inside the pods and see
 the details of containers (sidecars) running in a pod and then extended it so all
 containers could be viewed at once.

ice lists useful information about the (sidecar) containers present inside a
 pod, useful for trouble-shooting multi container issues. You can view volume, 
 image, port and executable configurations, along with current cpu and memory
 metrics all at the container level (requires metrics server)

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
T  CONTAINER    USED    REQUEST  LIMIT  %REQ    %LIMIT
I  app-init     0       0        0      -       -
S  app-watcher  0       1M       512M   -       -
S  app-broken   0       1M       512M   -       -
S  myapp        0.01Gi  1M       256M   550.09  2.15

```
### using labels
using labels you can search all pods that are part of a deployment where the label app matches demoprobe and list selected information about the containers in each pod, this example shows the currently configured probe information and gives details of configured startup, readiness and liveness probes of each container
``` shell
$ kubectl-ice probes -l app=demoprobe
PODNAME                      CONTAINER     PROBE     DELAY  PERIOD  TIMEOUT  SUCCESS  FAILURE  CHECK    ACTION
demo-probe-76b66d5766-ckhmt  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-ckhmt  web-frontend  liveness  5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-ckhmt  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/
demo-probe-76b66d5766-zk8pp  web-frontend  liveness  5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-zk8pp  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-zk8pp  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/

```
### Container status
most commands work the same way including the status command which also lets you see which container(s) are causing the restarts and by using the optional --previous flag you can view the containers previous exit code
``` shell
$ kubectl-ice status -l app=myapp --previous
T  PODNAME  CONTAINER    STATE       REASON  EXIT-CODE  SIGNAL  TIMESTAMP                      MESSAGE
S  web-pod  app-broken   Terminated  Error   1          0       2022-04-19 16:03:03 +0100 BST  -
S  web-pod  app-watcher  Terminated  Error   2          0       2022-04-19 15:59:18 +0100 BST  -
S  web-pod  myapp        -           -       -          -       -                              -
I  web-pod  app-init     -           -       -          -       -                              -

```
### Advanced labels
return memory requests size and limits of each container where the pods have an app label that matches useoddcpu and the container name is equal to nginx
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend
T  PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
S  demo-odd-cpu-5f947f9db4-4g2gk  web-frontend  2     1        1000   141.36    0.14
S  demo-odd-cpu-5f947f9db4-8pkpq  web-frontend  2     1        1000   132.28    0.13
S  demo-odd-cpu-5f947f9db4-9rdh4  web-frontend  2     1        1000   147.10    0.15
S  demo-odd-cpu-5f947f9db4-b227c  web-frontend  2     1        1000   157.30    0.16
S  demo-odd-cpu-5f947f9db4-dldv5  web-frontend  90    1        1000   8947.00   8.95
S  demo-odd-cpu-5f947f9db4-fg9wj  web-frontend  107   1        1000   10640.05  10.64
S  demo-odd-cpu-5f947f9db4-jnfck  web-frontend  2     1        1000   147.75    0.15
S  demo-odd-cpu-5f947f9db4-kktmw  web-frontend  2     1        1000   149.83    0.15
S  demo-odd-cpu-5f947f9db4-mgbpv  web-frontend  2     1        1000   138.21    0.14
S  demo-odd-cpu-5f947f9db4-p9mxv  web-frontend  2     1        1000   147.88    0.15
S  demo-odd-cpu-5f947f9db4-r5rk4  web-frontend  2     1        1000   140.73    0.14
S  demo-odd-cpu-5f947f9db4-rhktb  web-frontend  2     1        1000   145.02    0.15
S  demo-odd-cpu-5f947f9db4-tthwv  web-frontend  105   1        1000   10426.41  10.43
S  demo-odd-cpu-5f947f9db4-vqnm9  web-frontend  2     1        1000   131.28    0.13
S  demo-odd-cpu-5f947f9db4-w8w2z  web-frontend  2     1        1000   134.02    0.13
S  demo-odd-cpu-5f947f9db4-x65bq  web-frontend  2     1        1000   146.20    0.15

```
### Odditites and sorting
given the listed output above the optional --oddities flag picks out the containers that have a high cpu usage when compared to the other containers listed we also sort the list in descending order by the %REQ column
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend --oddities --sort '!%REQ'
T  PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
S  demo-odd-cpu-5f947f9db4-fg9wj  web-frontend  107   1        1000   10640.05  10.64
S  demo-odd-cpu-5f947f9db4-tthwv  web-frontend  105   1        1000   10426.41  10.43
S  demo-odd-cpu-5f947f9db4-dldv5  web-frontend  90    1        1000   8947.00   8.95

```
### Pod volumes
list all container volumes with mount points
``` shell
$ kubectl-ice volumes web-pod
CONTAINER    VOLUME                 TYPE       BACKING           SIZE  RO    MOUNT-POINT
app-init     kube-api-access-8z4tk  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-watcher  kube-api-access-8z4tk  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-broken   kube-api-access-8z4tk  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
myapp        app                    ConfigMap  app.py            -     false /myapp/
myapp        kube-api-access-8z4tk  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount

```


## Contributing

All feedback and contributions are welcome, if you want to raise an issue or help with fixes or features please [raise an issue to discuss](https://github.com/NimbleArchitect/kubectl-ice/issues)


## License
Licensed under Apache 2.0 see [LICENSE](https://github.com/NimbleArchitect/kubectl-pod/blob/main/LICENSE)
