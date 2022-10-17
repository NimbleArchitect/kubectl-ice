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
* Ability to read yaml from file or stdin for processing
* Can specify columns to output for a more custom view


[![asciicast](https://asciinema.org/a/512927.svg)](https://asciinema.org/a/512927)


## Contributing

If you like my work or find this program useful and want to say thanks you can reach me on twitter [@NimbleArchitect](https://twitter.com/nimblearchitect) or you can [Sponsor me](https://github.com/sponsors/NimbleArchitect) with github sponsors or [Buy Me A Coffee](https://www.buymeacoffee.com/NimbleArchitect)


All feedback and contributions are welcome, if you want to raise an issue or help with fixes or features please [raise an issue to discuss](https://github.com/NimbleArchitect/kubectl-ice/issues)


# Documentation
Full documentation can be found over at:

https://nimblearchitect.github.io/kubectl-ice

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
- download the required binary from the [downloads](https://nimblearchitect.github.io/kubectl-ice/dowloads/) page
- unzip and copy the kubectl-ice file to your path
- run ```kubectl-ice help``` to check its working

## Install from Source

clone and build the source using the following commands
```shell
git clone https://github.com/NimbleArchitect/kubectl-ice.git
cd kubectl-ice
make bin
```
then copy ./bin/kubectl-ice to somewhere in your path and run ```kubectl-ice version``` to check its working

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
Some example commands are listed below but full [usage instructions](https://nimblearchitect.github.io/kubectl-ice/documentation/#3_Usage) and [examples](https://nimblearchitect.github.io/kubectl-ice/documentation/#3.2_Example%20commands) can be found over at my website https://nimblearchitect.github.io/kubectl-ice/


### Single pod info
Shows the currently used memory along with the configured memory requests and limits of all containers (side cars) in the pod named web-pod
```
kubectl ice memory web-pod
```
### Named containers
the optional container flag (-c) searchs all selected pods and lists only containers that match the name web-frontend
```
kubectl ice command -c web-frontend
```

### Alternate status view
the tree flag shows the containers and pods in a tree view, with values calculated all the way up to the parent
```
kubectl ice status -l app=demoprobe --tree
```

### Excluding rows
use the --match flag to show only the output rows where the used memory column is greater than or equal to 3MB, this has the effect of exclusing any row where the used memory column is currently under 4096kB, the value 4096 can be replaced with any whole number in kilobytes
```
kubectl ice mem -l app=userandomcpu --match 'used>=4096'
```

### Extra selections
using the --select flag allows you to filter the pod selection to only pods that have a priorityClassName thats equal to system-cluster-critical, you can also match against priority
```
kubectl ice status --select 'priorityClassName=system-cluster-critical' -A
```

### Column labels
with the --node-label and --pod-label flags its possible to show the values of the labels as columns in the output table
```
kubectl ice status --node-label "beta.kubernetes.io/os" --pod-label "component" -n kube-system
```


## License
Licensed under Apache 2.0 see [LICENSE](https://github.com/NimbleArchitect/kubectl-pod/blob/main/LICENSE)
