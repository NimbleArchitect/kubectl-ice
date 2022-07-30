# kubectl-ice
A kubectl plugin that allows you to easily view advanced configuration of all containers
 that are running inside pods, I created it so I could peer inside the pods and see
 the details of the containers that are inside running pods and then extended it so all
 containers could be viewed at once.

ice lists detailed information about the containers present inside a
 pod, useful for trouble-shooting multi container issues. You can view volume, 
 image, port and executable configurations, along with current cpu and memory
 metrics all at the container level (requires metrics server)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nimblearchitect/kubectl-ice)
![GitHub](https://img.shields.io/github/license/NimbleArchitect/kubectl-ice)
![Github All Releases](https://img.shields.io/github/downloads/NimbleArchitect/kubectl-ice/total.svg?color=blue)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/NimbleArchitect/kubectl-ice/release)
![LGTM Alerts](https://img.shields.io/lgtm/alerts/github/NimbleArchitect/kubectl-ice)
[![CodeQL](https://github.com/NimbleArchitect/kubectl-ice/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/NimbleArchitect/kubectl-ice/actions/workflows/codeql-analysis.yml)

## Features:
* Runs on Windows, Linux and MacOS
* Only uses read permissions, no writes are called
* Tree view adds each container in a pod, then each pod in a replica or stateful set etc, all the way up to the node level
* Selectors work just like they do with the standard kubectl command
* Sortable output columns
* Can list all containers from all pods across all namespaces
* Include or exclude rows from output using the match flag, useful to exclude containers with low memory or cpu usage
* List only cpu and memory results that dont fall within range using the oddities flag
* Also displays information on init and ephemerial containers
* Pods can be filtered using their priority and priorityClassName
* Most sub commands utilize aliases meaning less typing (eg command and cmd are the same)
* Easily view securityContext details and POSIX capabilities
* Use the show-namespace flag to output the pods namespace


[![asciicast](https://asciinema.org/a/501737.svg)](https://asciinema.org/a/504766)

# Installation

## Install using krew

```
$ kubectl krew install ice
```
update with 
```
$ kubectl krew update
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
kubectl-ice capabilities  # Shows details of configured container POSIX capabilities
kubectl-ice command       # Retrieves the command line and any arguments specified at the container level
kubectl-ice cpu           # Show configured cpu size, limit and % usage of each container
kubectl-ice environment   # List the env name and value for each container
kubectl-ice help          # Help about any command
kubectl-ice image         # List the image name and pull status for each container
kubectl-ice ip            # List ip addresses of all pods in the namespace listed
kubectl-ice lifecycle     # Show lifecycle actions for each container in a named pod
kubectl-ice memory        # Show configured memory size, limit and % usage of each container
kubectl-ice ports         # Shows ports exposed by the containers in a pod
kubectl-ice probes        # Shows details of configured startup, readiness and liveness probes of each container
kubectl-ice restarts      # Show restart counts for each container in a named pod
kubectl-ice security      # Shows details of configured container security settings
kubectl-ice status        # List status of each container in a pod
kubectl-ice volumes       # Display container volumes and mount points
```

ice also supports all the standard kubectl flags in addition to:
```
Flags:
  -A, --all-namespaces                 List containers from pods in all namespaces
  -c, --container string               Container name. If set shows only the named containers
      --context string                 The name of the kubeconfig context to use
      --match string                   Filters out results, comma seperated list of COLUMN OP VALUE, where OP can be one of ==,<,>,<=,>= and != 
  -n, --namespace string               If present, the namespace scope for this CLI request
  -o, --output string                  Output format, currently csv, list, json and yaml are supported
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
  -l, --selector string                Selector (label query) to filter on
      --show-namespace                 Shows a column containing the pods namespace name for each container
  -t, --tree                           Display tree like view instead of the standard list
      --node-tree                      Displayes the tree with the nodes as the root
  -T  --show-type                      Show the container type column where:
                                            I = init container
                                            C = container
                                            E = ephemerial container
                                            P = Pod
                                            D = Deployment
                                            R = ReplicaSet
                                            A = DaemonSet
                                            S = StatefulSet
                                            N = Node

```
select subcommands also support the following flags
```
Flags:
  -d, --details          Display the timestamp instead of age along with the message column
  -p, --previous         Show previous state
  -r, --raw              Show raw uncooked values
      --sort string      Sort by column
      --oddities         Show only the outlier rows that dont fall within the computed range (requires min 5 rows in output)
```
all flags are optional, see usage instructions and examples for more info

## Examples
Some examples are listed below but full [usage instructions](https://github.com/NimbleArchitect/kubectl-pod/blob/main/docs/usage.md) and [examples](https://github.com/NimbleArchitect/kubectl-pod/blob/main/docs/examples.md) can be found in the [docs directory](https://github.com/NimbleArchitect/kubectl-pod/blob/main/docs)

### Single pod info
Shows the currently used memory along with the configured memory requests and limits of all containers (side cars) in the pod named web-pod
``` shell
$ kubectl-ice memory web-pod
CONTAINER    USED    REQUEST  LIMIT  %REQ    %LIMIT
app-init     0       0        0      -       -
app-watcher  5.25Mi  1M       512M   550.09  1.07
app-broken   0       1M       512M   -       -
myapp        5.23Mi  1M       256M   548.45  2.14

```
### Using labels
using labels you can search all pods that are part of a deployment where the label app matches demoprobe and list selected information about the containers in each pod, this example shows the currently configured probe information and gives details of configured startup, readiness and liveness probes of each container
``` shell
$ kubectl-ice probes -l app=demoprobe
PODNAME                      CONTAINER     PROBE     DELAY  PERIOD  TIMEOUT  SUCCESS  FAILURE  CHECK    ACTION
demo-probe-765fd4d8f7-7s5d7  web-frontend  liveness  10     5       1        1        3        Exec     /bin/true
demo-probe-765fd4d8f7-7s5d7  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-765fd4d8f7-7s5d7  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/
demo-probe-765fd4d8f7-cqr6m  web-frontend  liveness  10     5       1        1        3        Exec     /bin/true
demo-probe-765fd4d8f7-cqr6m  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-765fd4d8f7-cqr6m  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/

```
### Alternate status view
the tree flag shows the containers and pods in a tree view, with values calculated all the way up to the parent
``` shell
$ kubectl-ice status -l app=demoprobe --tree
NAMESPACE  NAME                                 READY  STARTED  RESTARTS  STATE    REASON  EXIT-CODE  SIGNAL  AGE
ice        Deployment/demo-probe                true   true     58        -        -       -          -       -
ice        └─ReplicaSet/demo-probe-765fd4d8f7   true   true     58        -        -       -          -       -
ice          └─Pod/demo-probe-765fd4d8f7-7s5d7  true   true     29        -        -       -          -       -
ice           └─Container/nginx                 true   true     0         Running  -       -          -       45h
ice           └─Container/web-frontend          true   true     29        Running  -       -          -       16m
ice          └─Pod/demo-probe-765fd4d8f7-cqr6m  true   true     29        -        -       -          -       -
ice           └─Container/nginx                 true   true     0         Running  -       -          -       45h
ice           └─Container/web-frontend          true   true     29        Running  -       -          -       16m

```
### Pick and un-mix
Using the -A flag to search all namespaces we can exclude all init containers with the --match T!=I flag. The -T flag is optional and is provided to show the init container type is not in the output
``` shell
$ kubectl-ice cpu -A -T --match T!=I
T  NAMESPACE    PODNAME                           CONTAINER                USED  REQUEST  LIMIT  %REQ      %LIMIT
C  default      web-pod                           app-watcher              5m    1m       1m     499.92    499.92
C  default      web-pod                           app-broken               0m    1m       1m     -         -
C  default      web-pod                           myapp                    5m    1m       1m     497.68    497.68
C  ice          demo-memory-7ddb58cd5b-dzp5x      web-frontend             993m  1m       1000m  99273.85  99.27
C  ice          demo-memory-7ddb58cd5b-dzp5x      nginx                    0m    1m       1000m  -         -
C  ice          demo-memory-7ddb58cd5b-hxkbt      web-frontend             997m  1m       1000m  99636.04  99.64
C  ice          demo-memory-7ddb58cd5b-hxkbt      nginx                    0m    1m       1000m  -         -
C  ice          demo-memory-7ddb58cd5b-wn7tt      web-frontend             997m  1m       1000m  99620.20  99.62
C  ice          demo-memory-7ddb58cd5b-wn7tt      nginx                    0m    1m       1000m  -         -
C  ice          demo-memory-7ddb58cd5b-xq2t4      web-frontend             994m  1m       1000m  99331.01  99.33
C  ice          demo-memory-7ddb58cd5b-xq2t4      nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-4clnc     web-frontend             129m  1m       1000m  12826.49  12.83
C  ice          demo-odd-cpu-5f947f9db4-4clnc     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-5z9w2     web-frontend             2m    1m       1000m  184.45    0.18
C  ice          demo-odd-cpu-5f947f9db4-5z9w2     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-62xjb     web-frontend             3m    1m       1000m  235.40    0.24
C  ice          demo-odd-cpu-5f947f9db4-62xjb     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-68f47     web-frontend             3m    1m       1000m  225.05    0.23
C  ice          demo-odd-cpu-5f947f9db4-68f47     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-7hlxl     web-frontend             3m    1m       1000m  223.38    0.22
C  ice          demo-odd-cpu-5f947f9db4-7hlxl     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-7r5s2     web-frontend             3m    1m       1000m  232.54    0.23
C  ice          demo-odd-cpu-5f947f9db4-7r5s2     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-8qpl5     web-frontend             3m    1m       1000m  222.98    0.22
C  ice          demo-odd-cpu-5f947f9db4-8qpl5     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-c5sv6     web-frontend             3m    1m       1000m  227.62    0.23
C  ice          demo-odd-cpu-5f947f9db4-c5sv6     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-c7scd     web-frontend             2m    1m       1000m  181.64    0.18
C  ice          demo-odd-cpu-5f947f9db4-c7scd     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-d6qz6     web-frontend             3m    1m       1000m  224.07    0.22
C  ice          demo-odd-cpu-5f947f9db4-d6qz6     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-jfwtd     web-frontend             2m    1m       1000m  196.39    0.20
C  ice          demo-odd-cpu-5f947f9db4-jfwtd     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-lkfs2     web-frontend             3m    1m       1000m  274.59    0.27
C  ice          demo-odd-cpu-5f947f9db4-lkfs2     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-q85pt     web-frontend             2m    1m       1000m  181.88    0.18
C  ice          demo-odd-cpu-5f947f9db4-q85pt     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-qdfnb     web-frontend             3m    1m       1000m  209.83    0.21
C  ice          demo-odd-cpu-5f947f9db4-qdfnb     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-tz4hj     web-frontend             3m    1m       1000m  236.00    0.24
C  ice          demo-odd-cpu-5f947f9db4-tz4hj     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-zsmwb     web-frontend             3m    1m       1000m  204.67    0.20
C  ice          demo-odd-cpu-5f947f9db4-zsmwb     nginx                    0m    1m       1000m  -         -
C  ice          demo-probe-765fd4d8f7-7s5d7       web-frontend             3m    125m     1000m  1.66      0.21
C  ice          demo-probe-765fd4d8f7-7s5d7       nginx                    0m    1m       1000m  -         -
C  ice          demo-probe-765fd4d8f7-cqr6m       web-frontend             3m    125m     1000m  2.27      0.28
C  ice          demo-probe-765fd4d8f7-cqr6m       nginx                    0m    1m       1000m  -         -
C  ice          demo-random-cpu-55954b64b4-2fnvf  web-frontend             148m  125m     1000m  118.33    14.79
C  ice          demo-random-cpu-55954b64b4-2fnvf  nginx                    0m    1m       1000m  -         -
C  ice          demo-random-cpu-55954b64b4-j5fgx  web-frontend             306m  125m     1000m  244.62    30.58
C  ice          demo-random-cpu-55954b64b4-j5fgx  nginx                    0m    1m       1000m  -         -
C  ice          demo-random-cpu-55954b64b4-kfc5q  web-frontend             260m  125m     1000m  207.35    25.92
C  ice          demo-random-cpu-55954b64b4-kfc5q  nginx                    0m    1m       1000m  -         -
C  ice          demo-random-cpu-55954b64b4-nd89h  web-frontend             521m  125m     1000m  416.55    52.07
C  ice          demo-random-cpu-55954b64b4-nd89h  nginx                    0m    1m       1000m  -         -
C  ice          web-pod                           app-watcher              5m    1m       1m     499.92    499.92
C  ice          web-pod                           app-broken               0m    1m       1m     -         -
C  ice          web-pod                           myapp                    5m    1m       1m     497.68    497.68
C  kube-system  coredns-78fcd69978-qnjtj          coredns                  1m    100m     0m     0.93      -
C  kube-system  etcd-minikube                     etcd                     9m    100m     0m     8.26      -
C  kube-system  kube-apiserver-minikube           kube-apiserver           34m   250m     0m     13.22     -
C  kube-system  kube-controller-manager-minikube  kube-controller-manager  7m    200m     0m     3.22      -
C  kube-system  kube-proxy-hdx8w                  kube-proxy               1m    0m       0m     -         -
C  kube-system  kube-scheduler-minikube           kube-scheduler           2m    100m     0m     1.37      -
C  kube-system  metrics-server-77c99ccb96-kjq9s   metrics-server           3m    100m     0m     2.40      -
C  kube-system  storage-provisioner               storage-provisioner      1m    0m       0m     -         -

```
### Filtered trees
The tree view also allows us to use the --match flag to filter based on resource type (T column) so we include deployments only providing us with a nice total of memory used for each deployment
``` shell
$ kubectl-ice mem -T --tree --match T==D
T  NAMESPACE  NAME                                      USED     REQUEST    LIMIT     %REQ     %LIMIT
D  ice        Deployment/demo-memory                    191.37Mi 11.44Mi    2334.59Mi 0.01     1.67
D  ice        Deployment/demo-odd-cpu                   119.67Mi 1556.40Mi  9338.38Mi 0.00     0.01
D  ice        Deployment/demo-probe                     6.43Mi   3.81Mi     976.56Mi  0.00     0.17
D  ice        Deployment/demo-random-cpu                42.19Mi  389.10Mi   2334.59Mi 0.00     0.01

```
### Status details
using the details flag displays the timestamp and message columns
``` shell
$ kubectl-ice status -l app=myapp --details
T  PODNAME  CONTAINER    READY  STARTED  RESTARTS  STATE       REASON            EXIT-CODE  SIGNAL  TIMESTAMP            MESSAGE
I  web-pod  app-init     true   -        0         Terminated  Completed         0          0       2022-07-28 18:16:48  -
C  web-pod  app-broken   false  false    191       Waiting     CrashLoopBackOff  -          -       -                    back-off 5m0s restarting failed
C  web-pod  app-watcher  true   true     0         Running     -                 -          -       2022-07-28 18:17:07  -
C  web-pod  myapp        true   true     0         Running     -                 -          -       2022-07-28 18:17:32  -

```
### Container status
most commands work the same way including the status command which also lets you see which container(s) are causing the restarts and by using the optional --previous flag you can view the containers previous exit code
``` shell
$ kubectl-ice status -l app=myapp --previous
PODNAME  CONTAINER    STATE       REASON  EXIT-CODE  SIGNAL  TIMESTAMP            MESSAGE
web-pod  app-init     -           -       -          -       -                    -
web-pod  app-broken   Terminated  Error   1          0       2022-07-30 15:57:50  -
web-pod  app-watcher  -           -       -          -       -                    -
web-pod  myapp        -           -       -          -       -                    -

```
### Advanced labels
return cpu requests size and limits of each container where the pods have an app label that matches useoddcpu and the container name is equal to web-frontend
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-4clnc  web-frontend  129m  1m       1000m  12826.49  12.83
demo-odd-cpu-5f947f9db4-5z9w2  web-frontend  2m    1m       1000m  184.45    0.18
demo-odd-cpu-5f947f9db4-62xjb  web-frontend  3m    1m       1000m  235.40    0.24
demo-odd-cpu-5f947f9db4-68f47  web-frontend  3m    1m       1000m  225.05    0.23
demo-odd-cpu-5f947f9db4-7hlxl  web-frontend  3m    1m       1000m  223.38    0.22
demo-odd-cpu-5f947f9db4-7r5s2  web-frontend  3m    1m       1000m  232.54    0.23
demo-odd-cpu-5f947f9db4-8qpl5  web-frontend  3m    1m       1000m  222.98    0.22
demo-odd-cpu-5f947f9db4-c5sv6  web-frontend  3m    1m       1000m  227.62    0.23
demo-odd-cpu-5f947f9db4-c7scd  web-frontend  2m    1m       1000m  181.64    0.18
demo-odd-cpu-5f947f9db4-d6qz6  web-frontend  3m    1m       1000m  224.07    0.22
demo-odd-cpu-5f947f9db4-jfwtd  web-frontend  2m    1m       1000m  196.39    0.20
demo-odd-cpu-5f947f9db4-lkfs2  web-frontend  3m    1m       1000m  274.59    0.27
demo-odd-cpu-5f947f9db4-q85pt  web-frontend  2m    1m       1000m  181.88    0.18
demo-odd-cpu-5f947f9db4-qdfnb  web-frontend  3m    1m       1000m  209.83    0.21
demo-odd-cpu-5f947f9db4-tz4hj  web-frontend  3m    1m       1000m  236.00    0.24
demo-odd-cpu-5f947f9db4-zsmwb  web-frontend  3m    1m       1000m  204.67    0.20

```
### Odditites and sorting
given the listed output above the optional --oddities flag picks out the containers that have a high cpu usage when compared to the other containers listed we also sort the list in descending order by the %REQ column
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend --oddities --sort '!%REQ'
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-4clnc  web-frontend  129m  1m       1000m  12826.49  12.83

```


## Contributing


If you like my work or find this program useful and want to say thanks you can reach me on twitter [@NimbleArchitect](https://twitter.com/nimblearchitect) or you can [Buy Me A Coffee](https://www.buymeacoffee.com/NimbleArchitect)


All feedback and contributions are welcome, if you want to raise an issue or help with fixes or features please [raise an issue to discuss](https://github.com/NimbleArchitect/kubectl-ice/issues)

## License
Licensed under Apache 2.0 see [LICENSE](https://github.com/NimbleArchitect/kubectl-pod/blob/main/LICENSE)
