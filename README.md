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
* List all the containers in a kubernetes pod including Init and Ephemeral containers
* Include or exclude rows from output using the match flag, useful to exclude containers with low memory or cpu usage
* List only cpu and memory results that dont fall within range using the oddities flag
* Also displays information on init and ephemerial containers
* Pods can be filtered using their priority and priorityClassName
* Most sub commands utilize aliases meaning less typing (eg command and cmd are the same)
* Easily view securityContext details and POSIX capabilities
* Use the show-namespace flag to output the pods namespace


[![asciicast](https://asciinema.org/a/512927.svg)](https://asciinema.org/a/512927)


## Contributing

If you like my work or find this program useful and want to say thanks you can reach me on twitter [@NimbleArchitect](https://twitter.com/nimblearchitect) or you can [Sponsor me](https://github.com/sponsors/NimbleArchitect) with github sponsors or [Buy Me A Coffee](https://www.buymeacoffee.com/NimbleArchitect)


All feedback and contributions are welcome, if you want to raise an issue or help with fixes or features please [raise an issue to discuss](https://github.com/NimbleArchitect/kubectl-ice/issues)


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
 if k is aliased to kubectl you can type ```k ice status``` instead of ```kubectl-ice status```


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
      --annotation string              Show the selected annotation as a column
  -c, --container string               Container name. If set shows only the named containers
      --context string                 The name of the kubeconfig context to use
  -m, --match string                   Filters out results, comma seperated list of COLUMN OP VALUE, where OP can be one of ==,<,>,<=,>= and != 
  -M, --match-only string              Filters out results but only calculates up visible rows
  -n, --namespace string               If present, the namespace scope for this CLI request
      --node-label string              Show the selected node label as a column
      --node-tree                      Displayes the tree with the nodes as the root
  -o, --output string                  Output format, currently csv, list, json and yaml are supported
      --pod-label string               Show the selected pod label as a column
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
  -l, --selector string                Selector (label query) to filter on
      --show-namespace                 Shows a column containing the pods namespace name for each container
  -t, --tree                           Display tree like view instead of the standard list
      --node-tree                      Displayes the tree with the nodes as the root
      --show-node                      Show the node name column
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
CONTAINER       USED    REQUEST  LIMIT  %REQ    %LIMIT
app-init        0       0        0      -       -
app-watcher     4.99Mi  1M       512M   523.06  1.02
app-broken      0       1M       512M   -       -
myapp           4.98Mi  1M       256M   521.83  2.04
debugger-k5znj  0       0        0      -       -

```
### Using labels
using labels you can search all pods that are part of a deployment where the label app matches demoprobe and list selected information about the containers in each pod, this example shows the currently configured probe information and gives details of configured startup, readiness and liveness probes of each container
``` shell
$ kubectl-ice probes -l app=demoprobe
PODNAME                      CONTAINER     PROBE     DELAY  PERIOD  TIMEOUT  SUCCESS  FAILURE  CHECK    ACTION
demo-probe-765fd4d8f7-n6kc7  web-frontend  liveness  10     5       1        1        3        Exec     /bin/true
demo-probe-765fd4d8f7-n6kc7  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-765fd4d8f7-n6kc7  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/
demo-probe-765fd4d8f7-x2zr6  web-frontend  liveness  10     5       1        1        3        Exec     /bin/true
demo-probe-765fd4d8f7-x2zr6  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-765fd4d8f7-x2zr6  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/

```
### Alternate status view
the tree flag shows the containers and pods in a tree view, with values calculated all the way up to the parent
``` shell
$ kubectl-ice status -l app=demoprobe --tree
NAMESPACE  NAME                                 READY  STARTED  RESTARTS  STATE    REASON  EXIT-CODE  SIGNAL  AGE
ice        Deployment/demo-probe                true   true     20        -        -       -          -       -
ice        └─ReplicaSet/demo-probe-765fd4d8f7   true   true     20        -        -       -          -       -
ice          └─Pod/demo-probe-765fd4d8f7-n6kc7  true   true     10        Running  -       -          -       19h
ice           └─Container/nginx                 true   true     0         Running  -       -          -       19h
ice           └─Container/web-frontend          true   true     10        Running  -       -          -       25m
ice          └─Pod/demo-probe-765fd4d8f7-x2zr6  true   true     10        Running  -       -          -       19h
ice           └─Container/nginx                 true   true     0         Running  -       -          -       19h
ice           └─Container/web-frontend          true   true     10        Running  -       -          -       25m

```
### Pick and un-mix
Using the -A flag to search all namespaces we can exclude all standard containers with the --match T!=C flag. The -T flag is optional and is provided to show that only Init and Ephemeral containers are displayed
``` shell
$ kubectl-ice cpu -A -T --match T!=C
T  NAMESPACE  PODNAME                           CONTAINER       USED  REQUEST  LIMIT  %REQ  %LIMIT
I  default    web-pod                           app-init        0m    0m       0m     -     -
I  ice        demo-memory-7ddb58cd5b-5kf9g      init-myservice  0m    1m       100m   -     -
I  ice        demo-memory-7ddb58cd5b-csbds      init-myservice  0m    1m       100m   -     -
I  ice        demo-memory-7ddb58cd5b-d4zwp      init-myservice  0m    1m       100m   -     -
I  ice        demo-memory-7ddb58cd5b-pdm9c      init-myservice  0m    1m       100m   -     -
I  ice        demo-odd-cpu-5f947f9db4-2g7p2     init-myservice  0m    100m     100m   -     -
I  ice        demo-odd-cpu-5f947f9db4-59jm6     init-myservice  0m    100m     100m   -     -
I  ice        demo-odd-cpu-5f947f9db4-6gzw7     init-myservice  0m    100m     100m   -     -
I  ice        demo-odd-cpu-5f947f9db4-6s97l     init-myservice  0m    100m     100m   -     -
E  ice        demo-odd-cpu-5f947f9db4-6s97l     debugger-vvz4z  0m    0m       0m     -     -
I  ice        demo-odd-cpu-5f947f9db4-86mb8     init-myservice  0m    100m     100m   -     -
I  ice        demo-odd-cpu-5f947f9db4-cwvdq     init-myservice  0m    100m     100m   -     -
I  ice        demo-odd-cpu-5f947f9db4-dcg8p     init-myservice  0m    100m     100m   -     -
I  ice        demo-odd-cpu-5f947f9db4-fhs8q     init-myservice  0m    100m     100m   -     -
I  ice        demo-odd-cpu-5f947f9db4-gzcrm     init-myservice  0m    100m     100m   -     -
I  ice        demo-odd-cpu-5f947f9db4-hf872     init-myservice  0m    100m     100m   -     -
I  ice        demo-odd-cpu-5f947f9db4-hft68     init-myservice  0m    100m     100m   -     -
I  ice        demo-odd-cpu-5f947f9db4-jp8fw     init-myservice  0m    100m     100m   -     -
I  ice        demo-odd-cpu-5f947f9db4-k2gtp     init-myservice  0m    100m     100m   -     -
I  ice        demo-odd-cpu-5f947f9db4-kj8s7     init-myservice  0m    100m     100m   -     -
I  ice        demo-odd-cpu-5f947f9db4-qtxp2     init-myservice  0m    100m     100m   -     -
I  ice        demo-odd-cpu-5f947f9db4-vg2d5     init-myservice  0m    100m     100m   -     -
I  ice        demo-random-cpu-55954b64b4-9t7m2  init-myservice  0m    120m     120m   -     -
I  ice        demo-random-cpu-55954b64b4-km6bg  init-myservice  0m    120m     120m   -     -
I  ice        demo-random-cpu-55954b64b4-knc6n  init-myservice  0m    120m     120m   -     -
I  ice        demo-random-cpu-55954b64b4-vr4hg  init-myservice  0m    120m     120m   -     -
I  ice        web-pod                           app-init        0m    0m       0m     -     -
E  ice        web-pod                           debugger-k5znj  0m    0m       0m     -     -

```
### Filtered trees
The tree view also allows us to use the --match flag to filter based on resource type (T column) so we include deployments only providing us with a nice total of memory used for each deployment
``` shell
$ kubectl-ice mem -T --tree --match T==D
T  NAMESPACE  NAME                        USED      REQUEST    LIMIT      %REQ  %LIMIT
D  ice        Deployment/demo-memory      328.13Mi  11.44Mi    2334.59Mi  0.01  2.87
D  ice        Deployment/demo-odd-cpu     103.47Mi  1459.12Mi  8754.73Mi  0.00  0.01
D  ice        Deployment/demo-random-cpu  15.69Mi   194.55Mi   1167.30Mi  0.00  0.01
D  ice        Deployment/demo-odd-cpu     7.14Mi    97.27Mi    583.65Mi   0.00  0.01
D  ice        Deployment/demo-probe       6.42Mi    3.81Mi     976.56Mi   0.00  0.17
D  ice        Deployment/demo-random-cpu  21.33Mi   194.55Mi   1167.30Mi  0.00  0.01

```
### Status details
using the details flag displays the timestamp and message columns
``` shell
$ kubectl-ice status -l app=myapp --details
T  PODNAME  CONTAINER       READY  STARTED  RESTARTS  STATE       REASON            EXIT-CODE  SIGNAL  TIMESTAMP            MESSAGE
I  web-pod  app-init        true   -        0         Terminated  Completed         0          0       2022-08-04 19:00:53  -
C  web-pod  app-broken      false  false    68        Waiting     CrashLoopBackOff  -          -       -                    back-off 5m0s restarting failed
C  web-pod  app-watcher     true   true     0         Running     -                 -          -       2022-08-04 19:00:59  -
C  web-pod  myapp           true   true     0         Running     -                 -          -       2022-08-04 19:01:12  -
E  web-pod  debugger-k5znj  false  -        0         Terminated  Completed         0          0       2022-08-04 19:02:59  -

```
### Container status
most commands work the same way including the status command which also lets you see which container(s) are causing the restarts and by using the optional --previous flag you can view the containers previous exit code
``` shell
$ kubectl-ice status -l app=myapp --previous
PODNAME  CONTAINER       STATE       REASON  EXIT-CODE  SIGNAL  TIMESTAMP            MESSAGE
web-pod  app-init        -           -       -          -       -                    -
web-pod  app-broken      Terminated  Error   1          0       2022-08-05 13:59:58  -
web-pod  app-watcher     -           -       -          -       -                    -
web-pod  myapp           -           -       -          -       -                    -
web-pod  debugger-k5znj  -           -       -          -       -                    -

```
### Advanced labels
return cpu requests size and limits of each container where the pods have an app label that matches useoddcpu and the container name is equal to web-frontend
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-2g7p2  web-frontend  3m    1m       1000m  239.26    0.24
demo-odd-cpu-5f947f9db4-59jm6  web-frontend  3m    1m       1000m  228.24    0.23
demo-odd-cpu-5f947f9db4-6gzw7  web-frontend  137m  1m       1000m  13605.66  13.61
demo-odd-cpu-5f947f9db4-6s97l  web-frontend  3m    1m       1000m  242.12    0.24
demo-odd-cpu-5f947f9db4-86mb8  web-frontend  3m    1m       1000m  236.35    0.24
demo-odd-cpu-5f947f9db4-cwvdq  web-frontend  135m  1m       1000m  13408.25  13.41
demo-odd-cpu-5f947f9db4-dcg8p  web-frontend  2m    1m       1000m  189.21    0.19
demo-odd-cpu-5f947f9db4-fhs8q  web-frontend  3m    1m       1000m  231.62    0.23
demo-odd-cpu-5f947f9db4-gzcrm  web-frontend  3m    1m       1000m  239.25    0.24
demo-odd-cpu-5f947f9db4-hf872  web-frontend  2m    1m       1000m  196.35    0.20
demo-odd-cpu-5f947f9db4-hft68  web-frontend  3m    1m       1000m  219.94    0.22
demo-odd-cpu-5f947f9db4-jp8fw  web-frontend  3m    1m       1000m  235.47    0.24
demo-odd-cpu-5f947f9db4-k2gtp  web-frontend  135m  1m       1000m  13417.88  13.42
demo-odd-cpu-5f947f9db4-kj8s7  web-frontend  3m    1m       1000m  229.13    0.23
demo-odd-cpu-5f947f9db4-qtxp2  web-frontend  3m    1m       1000m  262.48    0.26
demo-odd-cpu-5f947f9db4-vg2d5  web-frontend  3m    1m       1000m  252.21    0.25

```
### Odditites and sorting
given the listed output above the optional --oddities flag picks out the containers that have a high cpu usage when compared to the other containers listed we also sort the list in descending order by the %REQ column
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend --oddities --sort '!%REQ'
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-6gzw7  web-frontend  137m  1m       1000m  13605.66  13.61
demo-odd-cpu-5f947f9db4-k2gtp  web-frontend  135m  1m       1000m  13417.88  13.42
demo-odd-cpu-5f947f9db4-cwvdq  web-frontend  135m  1m       1000m  13408.25  13.41

```


## License
Licensed under Apache 2.0 see [LICENSE](https://github.com/NimbleArchitect/kubectl-pod/blob/main/LICENSE)
