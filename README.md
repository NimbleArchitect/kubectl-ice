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
* Selectors work just like they do with the standard kubectl command
* Sortable output columns
* Can list all containers from all pods across all namespaces
* Exclude rows from output using the match flag, useful to exclude containers with low memory or cpu usage
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
  -T  --show-type                      Show the container type column where I = init container, S = standard container and E = ephemerial container
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
app-watcher  6.52Mi  1M       512M   683.21  1.33
app-broken   0       1M       512M   -       -
myapp        6.79Mi  1M       256M   711.48  2.78

```
### Using labels
using labels you can search all pods that are part of a deployment where the label app matches demoprobe and list selected information about the containers in each pod, this example shows the currently configured probe information and gives details of configured startup, readiness and liveness probes of each container
``` shell
$ kubectl-ice probes -l app=demoprobe
PODNAME                      CONTAINER     PROBE     DELAY  PERIOD  TIMEOUT  SUCCESS  FAILURE  CHECK    ACTION
demo-probe-765fd4d8f7-r5rq7  web-frontend  liveness  10     5       1        1        3        Exec     /bin/true
demo-probe-765fd4d8f7-r5rq7  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-765fd4d8f7-r5rq7  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/
demo-probe-765fd4d8f7-wxt7f  web-frontend  liveness  10     5       1        1        3        Exec     /bin/true
demo-probe-765fd4d8f7-wxt7f  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-765fd4d8f7-wxt7f  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/

```
### Alternate status view
the tree flag shows the containers and pods in a tree like view
``` shell
$ kubectl-ice status -l app=myapp --tree
NAMESPACE  NAME                      READY  STARTED  RESTARTS  STATE       REASON            EXIT-CODE  SIGNAL  AGE
ice        Pod/web-pod               -      -        0         Running     -                 -          -       104m
ice        └─InitContainer/app-init  true   -        0         Terminated  Completed         0          0       104m
ice        └─Container/app-broken    false  false    24        Waiting     CrashLoopBackOff  -          -       -
ice        └─Container/app-watcher   true   true     0         Running     -                 -          -       103m
ice        └─Container/myapp         true   true     0         Running     -                 -          -       103m

```
### Pick and un-mix
Using the -A flag to search all namespaces we can exclude all init containers with the --match T!=I flag. The -T flag is optional and is provided to show the init container type is not in the output
``` shell
$ kubectl-ice cpu -A -T --match T!=I
T  NAMESPACE    PODNAME                           CONTAINER                USED  REQUEST  LIMIT  %REQ      %LIMIT
S  default      web-pod                           app-watcher              5m    1m       1m     495.46    495.46
S  default      web-pod                           app-broken               0m    1m       1m     -         -
S  default      web-pod                           myapp                    6m    1m       1m     516.79    516.79
S  ice          demo-memory-7ddb58cd5b-g76gd      web-frontend             997m  1m       1000m  99633.77  99.63
S  ice          demo-memory-7ddb58cd5b-g76gd      nginx                    0m    1m       1000m  -         -
S  ice          demo-memory-7ddb58cd5b-ngbng      web-frontend             996m  1m       1000m  99528.84  99.53
S  ice          demo-memory-7ddb58cd5b-ngbng      nginx                    0m    1m       1000m  -         -
S  ice          demo-memory-7ddb58cd5b-r22kp      web-frontend             997m  1m       1000m  99671.11  99.67
S  ice          demo-memory-7ddb58cd5b-r22kp      nginx                    0m    1m       1000m  -         -
S  ice          demo-memory-7ddb58cd5b-r6bgn      web-frontend             997m  1m       1000m  99683.65  99.68
S  ice          demo-memory-7ddb58cd5b-r6bgn      nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-2wgdp     web-frontend             3m    1m       1000m  276.45    0.28
S  ice          demo-odd-cpu-5f947f9db4-2wgdp     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-56hnq     web-frontend             2m    1m       1000m  199.89    0.20
S  ice          demo-odd-cpu-5f947f9db4-56hnq     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-77q9n     web-frontend             123m  1m       1000m  12287.04  12.29
S  ice          demo-odd-cpu-5f947f9db4-77q9n     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-9kwbm     web-frontend             3m    1m       1000m  234.97    0.23
S  ice          demo-odd-cpu-5f947f9db4-9kwbm     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-9t7vg     web-frontend             3m    1m       1000m  209.26    0.21
S  ice          demo-odd-cpu-5f947f9db4-9t7vg     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-g2258     web-frontend             3m    1m       1000m  271.59    0.27
S  ice          demo-odd-cpu-5f947f9db4-g2258     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-hd6ld     web-frontend             3m    1m       1000m  204.26    0.20
S  ice          demo-odd-cpu-5f947f9db4-hd6ld     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-mnznh     web-frontend             113m  1m       1000m  11287.89  11.29
S  ice          demo-odd-cpu-5f947f9db4-mnznh     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-n9zvc     web-frontend             3m    1m       1000m  215.90    0.22
S  ice          demo-odd-cpu-5f947f9db4-n9zvc     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-qtxmf     web-frontend             3m    1m       1000m  238.15    0.24
S  ice          demo-odd-cpu-5f947f9db4-qtxmf     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-ssv2x     web-frontend             3m    1m       1000m  253.31    0.25
S  ice          demo-odd-cpu-5f947f9db4-ssv2x     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-v24n4     web-frontend             3m    1m       1000m  261.50    0.26
S  ice          demo-odd-cpu-5f947f9db4-v24n4     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-wfg78     web-frontend             3m    1m       1000m  240.82    0.24
S  ice          demo-odd-cpu-5f947f9db4-wfg78     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-x5k44     web-frontend             2m    1m       1000m  198.79    0.20
S  ice          demo-odd-cpu-5f947f9db4-x5k44     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-xbrhk     web-frontend             3m    1m       1000m  254.32    0.25
S  ice          demo-odd-cpu-5f947f9db4-xbrhk     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-zt5vw     web-frontend             3m    1m       1000m  229.52    0.23
S  ice          demo-odd-cpu-5f947f9db4-zt5vw     nginx                    0m    1m       1000m  -         -
S  ice          demo-probe-765fd4d8f7-r5rq7       web-frontend             3m    125m     1000m  2.02      0.25
S  ice          demo-probe-765fd4d8f7-r5rq7       nginx                    1m    1m       1000m  1.24      0.00
S  ice          demo-probe-765fd4d8f7-wxt7f       web-frontend             3m    125m     1000m  2.05      0.26
S  ice          demo-probe-765fd4d8f7-wxt7f       nginx                    0m    1m       1000m  -         -
S  ice          demo-random-cpu-55954b64b4-l4b44  web-frontend             145m  125m     1000m  115.28    14.41
S  ice          demo-random-cpu-55954b64b4-l4b44  nginx                    0m    1m       1000m  -         -
S  ice          demo-random-cpu-55954b64b4-nwm9c  web-frontend             249m  125m     1000m  199.01    24.88
S  ice          demo-random-cpu-55954b64b4-nwm9c  nginx                    1m    1m       1000m  1.73      0.00
S  ice          demo-random-cpu-55954b64b4-vnrk5  web-frontend             279m  125m     1000m  222.59    27.82
S  ice          demo-random-cpu-55954b64b4-vnrk5  nginx                    1m    1m       1000m  0.77      0.00
S  ice          demo-random-cpu-55954b64b4-xnxhd  web-frontend             254m  125m     1000m  202.81    25.35
S  ice          demo-random-cpu-55954b64b4-xnxhd  nginx                    1m    1m       1000m  1.10      0.00
S  ice          web-pod                           app-watcher              5m    1m       1m     495.46    495.46
S  ice          web-pod                           app-broken               0m    1m       1m     -         -
S  ice          web-pod                           myapp                    6m    1m       1m     516.79    516.79
S  kube-system  coredns-78fcd69978-qnjtj          coredns                  1m    100m     0m     0.88      -
S  kube-system  etcd-minikube                     etcd                     8m    100m     0m     7.16      -
S  kube-system  kube-apiserver-minikube           kube-apiserver           29m   250m     0m     11.42     -
S  kube-system  kube-controller-manager-minikube  kube-controller-manager  11m   200m     0m     5.29      -
S  kube-system  kube-proxy-hdx8w                  kube-proxy               1m    0m       0m     -         -
S  kube-system  kube-scheduler-minikube           kube-scheduler           2m    100m     0m     1.06      -
S  kube-system  metrics-server-77c99ccb96-kjq9s   metrics-server           2m    100m     0m     1.78      -
S  kube-system  storage-provisioner               storage-provisioner      1m    0m       0m     -         -

```
### Status details
using the details flag displays the timestamp and message columns
``` shell
$ kubectl-ice status -l app=myapp --details
T  PODNAME  CONTAINER    READY  STARTED  RESTARTS  STATE       REASON            EXIT-CODE  SIGNAL  TIMESTAMP            MESSAGE
I  web-pod  app-init     true   -        0         Terminated  Completed         0          0       2022-07-13 16:24:31  -
S  web-pod  app-broken   false  false    24        Waiting     CrashLoopBackOff  -          -       -                    back-off 5m0s restarting failed
S  web-pod  app-watcher  true   true     0         Running     -                 -          -       2022-07-13 16:25:04  -
S  web-pod  myapp        true   true     0         Running     -                 -          -       2022-07-13 16:25:17  -

```
### Container status
most commands work the same way including the status command which also lets you see which container(s) are causing the restarts and by using the optional --previous flag you can view the containers previous exit code
``` shell
$ kubectl-ice status -l app=myapp --previous
PODNAME  CONTAINER    STATE       REASON  EXIT-CODE  SIGNAL  TIMESTAMP            MESSAGE
web-pod  app-init     -           -       -          -       -                    -
web-pod  app-broken   Terminated  Error   1          0       2022-07-13 18:07:49  -
web-pod  app-watcher  -           -       -          -       -                    -
web-pod  myapp        -           -       -          -       -                    -

```
### Advanced labels
return memory requests size and limits of each container where the pods have an app label that matches useoddcpu and the container name is equal to web-frontend
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-2wgdp  web-frontend  3m    1m       1000m  276.45    0.28
demo-odd-cpu-5f947f9db4-56hnq  web-frontend  2m    1m       1000m  199.89    0.20
demo-odd-cpu-5f947f9db4-77q9n  web-frontend  123m  1m       1000m  12287.04  12.29
demo-odd-cpu-5f947f9db4-9kwbm  web-frontend  3m    1m       1000m  234.97    0.23
demo-odd-cpu-5f947f9db4-9t7vg  web-frontend  3m    1m       1000m  209.26    0.21
demo-odd-cpu-5f947f9db4-g2258  web-frontend  3m    1m       1000m  271.59    0.27
demo-odd-cpu-5f947f9db4-hd6ld  web-frontend  3m    1m       1000m  204.26    0.20
demo-odd-cpu-5f947f9db4-mnznh  web-frontend  113m  1m       1000m  11287.89  11.29
demo-odd-cpu-5f947f9db4-n9zvc  web-frontend  3m    1m       1000m  215.90    0.22
demo-odd-cpu-5f947f9db4-qtxmf  web-frontend  3m    1m       1000m  238.15    0.24
demo-odd-cpu-5f947f9db4-ssv2x  web-frontend  3m    1m       1000m  253.31    0.25
demo-odd-cpu-5f947f9db4-v24n4  web-frontend  3m    1m       1000m  261.50    0.26
demo-odd-cpu-5f947f9db4-wfg78  web-frontend  3m    1m       1000m  240.82    0.24
demo-odd-cpu-5f947f9db4-x5k44  web-frontend  2m    1m       1000m  198.79    0.20
demo-odd-cpu-5f947f9db4-xbrhk  web-frontend  3m    1m       1000m  254.32    0.25
demo-odd-cpu-5f947f9db4-zt5vw  web-frontend  3m    1m       1000m  229.52    0.23

```
### Odditites and sorting
given the listed output above the optional --oddities flag picks out the containers that have a high cpu usage when compared to the other containers listed we also sort the list in descending order by the %REQ column
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend --oddities --sort '!%REQ'
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-77q9n  web-frontend  123m  1m       1000m  12287.04  12.29
demo-odd-cpu-5f947f9db4-mnznh  web-frontend  113m  1m       1000m  11287.89  11.29

```
### Pod volumes
list all container volumes with mount points
``` shell
$ kubectl-ice volumes web-pod
CONTAINER    VOLUME                 TYPE       BACKING           SIZE  RO    MOUNT-POINT
app-init     kube-api-access-7cx8b  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-watcher  app                    ConfigMap  app.py            -     false /myapp/
app-watcher  kube-api-access-7cx8b  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-broken   kube-api-access-7cx8b  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
myapp        app                    ConfigMap  app.py            -     false /myapp/
myapp        kube-api-access-7cx8b  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount

```


## Contributing


If you like my work or find this program useful and want to say thanks you can reach me on twitter [@NimbleArchitect](https://twitter.com/nimblearchitect) or you can [Buy Me A Coffee](https://www.buymeacoffee.com/NimbleArchitect)


All feedback and contributions are welcome, if you want to raise an issue or help with fixes or features please [raise an issue to discuss](https://github.com/NimbleArchitect/kubectl-ice/issues)

## License
Licensed under Apache 2.0 see [LICENSE](https://github.com/NimbleArchitect/kubectl-pod/blob/main/LICENSE)
