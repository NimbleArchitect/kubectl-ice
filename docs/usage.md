# Contents
* [Introduction](#introduction)
* [Usage](#usage)
* [Flags](#flags)
* [Capabilities](#capabilities)
  * [Examples](#example)
* [Command](#command)
  * [Examples](#example-1)
* [CPU](#cpu)
  * [Examples](#example-2)
* [Environment](#environment)
  * [Examples](#example-3)
* [Image](#image)
  * [Examples](#example-4)
* [IP](#ip)
  * [Examples](#example-5)
* [Memory](#memory)
  * [Examples](#example-6)
* [Ports](#ports)
  * [Examples](#example-7)
* [Probes](#probes)
  * [Examples](#example-8)
* [Restarts](#restarts)
  * [Examples](#example-9)
* [Security](#security)
  * [Examples](#example-10)
* [Status](#status)
  * [Examples](#example-11)
* [Volumes](#volumes)
  * [Examples](#example-12)

## Introduction
A kubectl plugin that lets you can see the running configuration of all containers
 that are running inside pods, I created it so I could peer inside the pods and see
 the details of containers (sidecars) running in a pod and then extended it so all
 containers could be viewed at once.   

ice lists useful information about the sidecar containers present inside a
 pod, useful for trouble shooting multi container issues you can view volume, 
 image, port and executable configurations, along with current cpu and memory
  metrics all at the container level (requires metrics server)

## Usage
ice usage is split in to sub commands with each following commands are available for `kubectl ice`
```
kubectl-ice capabilities  # Shows details of configured container capabilities
kubectl-ice command       # Retrieves the command line and any arguments specified at the container level
kubectl-ice cpu           # Show configured cpu size, limit and % usage of each container
kubectl-ice environment   # List the env name and value for each container
kubectl-ice help          # Help about any command
kubectl-ice image         # List the image name and pull status for each container
kubectl-ice ip            # List ip addresses of all pods in the namespace listed
kubectl-ice memory        # Show configured memory size, limit and % usage of each container
kubectl-ice ports         # Shows ports exposed by the containers in a pod
kubectl-ice probes        # Shows details of configured startup, readiness and liveness probes of each container
kubectl-ice restarts      # Show restart counts for each container in a named pod
kubectl-ice security      # Shows details of configured container security settings
kubectl-ice status        # List status of each container in a pod
kubectl-ice volumes       # Display container volumes and mount points
```

## Flags
All standard kubectl flags are supported including the beow, see the examples section of each command for usage info:
```
  -A, --all-namespaces                 List containers from pods in all namespaces
  -c, --container string               Container name. If set shows only the named containers containers in the pod
      --context string                 The name of the kubeconfig context to use
      --match string                   excludes results, comma seperated list of COLUMN OP VALUE, where OP can be one of ==,<,>,<=,>= and != 
  -n, --namespace string               If present, the namespace scope for this CLI request
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
  -l, --selector string                Selector (label query) to filter on
```
selected subcommands also support the following flags
```
  -p, --previous         show previous state
  -r, --raw              show raw uncooked values
      --sort string      Sort by column
      --oddities         show only the outlier rows that dont fall within the computed range (requires min 5 rows in output)
```

## Capabilities

``` shell
Usage:
  kubectl-ice capabilities [flags]

Aliases:
  capabilities, cap

Flags:
      --node-label string              Show the selected node label as a column
  -o, --output string                  Output format, currently csv, list, json and yaml are supported
      --pod-label string               Show the selected pod label as a column
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
      --show-namespace                 Show the namespace column
      --show-node                      Show the node name column
  -T, --show-type                      Show the container type column
  -t, --tree                           Display tree like view instead of the standard list


More information at: https://www.github.com/NimbleArchitect/kubectl-ice

```
also includes standard common kubectl flags

### Examples
``` shell
# List container capabilities from pods
kubectl-ice capabilities

# List container capabilities info from pods output in JSON format
kubectl-ice capabilities -o json

# List container capabilities info from a single pod
kubectl-ice capabilities my-pod-4jh36

# List capabilities info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice capabilities -c web-container

# List capabilities info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice capabilities -c web-container --sort '!CONTAINER'

# List capabilities info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice capabilities -c web-container --sort PODNAME

# List container capabilities info from all pods where label app matches web
kubectl-ice capabilities -l app=web

# List container capabilities info from all pods where the pod label app is either web or mail
kubectl-ice capabilities -l "app in (web,mail)"
```
## Command

``` shell
Usage:
  kubectl-ice command [flags]

Aliases:
  command, cmd, exec, args

Flags:
      --node-label string              Show the selected node label as a column
  -o, --output string                  Output format, currently csv, list, json and yaml are supported
      --pod-label string               Show the selected pod label as a column
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
      --show-namespace                 Show the namespace column
      --show-node                      Show the node name column
  -T, --show-type                      Show the container type column
  -t, --tree                           Display tree like view instead of the standard list


More information at: https://www.github.com/NimbleArchitect/kubectl-ice

```
also includes standard common kubectl flags

### Examples
``` shell
# List containers command info from pods
kubectl-ice command

# List container command info from pods output in JSON format
kubectl-ice command -o json

# List container command info from a single pod
kubectl-ice command my-pod-4jh36

# List command info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice command -c web-container

# List command info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice command -c web-container --sort '!CONTAINER'

# List command info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice command -c web-container --sort PODNAME

# List container command info from all pods where label app matches web
kubectl-ice command -l app=web

# List container command info from all pods where the pod label app is either web or mail
kubectl-ice command -l "app in (web,mail)"
```
## CPU

``` shell
Usage:
  kubectl-ice cpu [flags]

Flags:
  -i, --include-init                   include init container(s) in the output, by default init containers are hidden
      --node-label string              Show the selected node label as a column
      --oddities                       show only the outlier rows that dont fall within the computed range
  -o, --output string                  Output format, currently csv, list, json and yaml are supported
      --pod-label string               Show the selected pod label as a column
  -r, --raw                            show raw values
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
      --show-namespace                 Show the namespace column
      --show-node                      Show the node name column
  -T, --show-type                      Show the container type column
  -t, --tree                           Display tree like view instead of the standard list


More information at: https://www.github.com/NimbleArchitect/kubectl-ice

```
also includes standard common kubectl flags

### Examples
``` shell
# List containers cpu info from pods
kubectl-ice cpu

# List container cpu info from pods output in JSON format
kubectl-ice cpu -o json

# List container cpu info from a single pod
kubectl-ice cpu my-pod-4jh36

# List cpu info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice cpu -c web-container

# List cpu info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice cpu -c web-container --sort '!CONTAINER'

# List cpu info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice cpu -c web-container --sort PODNAME

# List container cpu info from all pods where label app matches web
kubectl-ice cpu -l app=web

# List container cpu info from all pods where the pod label app is either web or mail
kubectl-ice cpu -l "app in (web,mail)"
```
## Environment

``` shell
Usage:
  kubectl-ice environment [flags]

Aliases:
  environment, env, vars

Flags:
      --node-label string              Show the selected node label as a column
  -o, --output string                  Output format, currently csv, list, json and yaml are supported
      --pod-label string               Show the selected pod label as a column
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
      --show-namespace                 Show the namespace column
      --show-node                      Show the node name column
  -T, --show-type                      Show the container type column
      --translate                      read the configmap show its values
  -t, --tree                           Display tree like view instead of the standard list


More information at: https://www.github.com/NimbleArchitect/kubectl-ice

```
also includes standard common kubectl flags

### Examples
``` shell
# List containers env info from pods
kubectl-ice env

# List container env info from pods output in JSON format
kubectl-ice env -o json

# List container env info from a single pod
kubectl-ice env my-pod-4jh36

# List env info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice env -c web-container

# List env info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice env -c web-container --sort '!CONTAINER'

# List env info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice env -c web-container --sort PODNAME

# List container env info from all pods where label app matches web
kubectl-ice env -l app=web

# List container env info from all pods where the pod label app is either web or mail
kubectl-ice env -l "app in (web,mail)"
```
## Image

``` shell
Usage:
  kubectl-ice image [flags]

Aliases:
  image, im

Flags:
      --node-label string              Show the selected node label as a column
  -o, --output string                  Output format, currently csv, list, json and yaml are supported
      --pod-label string               Show the selected pod label as a column
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
      --show-namespace                 Show the namespace column
      --show-node                      Show the node name column
  -T, --show-type                      Show the container type column
  -t, --tree                           Display tree like view instead of the standard list


More information at: https://www.github.com/NimbleArchitect/kubectl-ice

```
also includes standard common kubectl flags

### Examples
``` shell
# List containers image info from pods
kubectl-ice image

# List container image info from pods output in JSON format
kubectl-ice image -o json

# List container image info from a single pod
kubectl-ice image my-pod-4jh36

# List image info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice image -c web-container

# List image info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice image -c web-container --sort '!CONTAINER'

# List image info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice image -c web-container --sort PODNAME

# List container image info from all pods where label app matches web
kubectl-ice image -l app=web

# List container image info from all pods where the pod label app is either web or mail
kubectl-ice image -l "app in (web,mail)"
```
## IP

``` shell
Usage:
  kubectl-ice ip [flags]

Flags:
      --node-label string              Show the selected node label as a column
  -o, --output string                  Output format, currently csv, list, json and yaml are supported
      --pod-label string               Show the selected pod label as a column
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
      --show-namespace                 Show the namespace column
      --show-node                      Show the node name column
  -T, --show-type                      Show the container type column
  -t, --tree                           Display tree like view instead of the standard list


More information at: https://www.github.com/NimbleArchitect/kubectl-ice

```
also includes standard common kubectl flags

### Examples
``` shell
# List IP address of pods
kubectl-ice ip

# List IP address of pods output in JSON format
kubectl-ice ip -o json

# List IP address a single pod
kubectl-ice ip my-pod-4jh36

# List IP address of all pods where label app matches web
kubectl-ice ip -l app=web

# List IP address of all pods where the pod label app is either web or mail
kubectl-ice ip -l "app in (web,mail)"
```
## Memory

``` shell
Usage:
  kubectl-ice memory [flags]

Aliases:
  memory, mem

Flags:
  -i, --include-init                   include init container(s) in the output, by default init containers are hidden
      --node-label string              Show the selected node label as a column
      --oddities                       show only the outlier rows that dont fall within the computed range
  -o, --output string                  Output format, currently csv, list, json and yaml are supported
      --pod-label string               Show the selected pod label as a column
  -r, --raw                            show raw values
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
      --show-namespace                 Show the namespace column
      --show-node                      Show the node name column
  -T, --show-type                      Show the container type column
      --size string                    allows conversion to the selected size rather then the default megabyte output (default "Mi")
  -t, --tree                           Display tree like view instead of the standard list


More information at: https://www.github.com/NimbleArchitect/kubectl-ice

```
also includes standard common kubectl flags

### Examples
``` shell
# List containers memory info from pods
kubectl-ice memory

# List container memory info from pods output in JSON format
kubectl-ice memory -o json

# List container memory info from a single pod
kubectl-ice memory my-pod-4jh36

# List memory info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice memory -c web-container

# List memory info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice memory -c web-container --sort '!CONTAINER'

# List memory info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice memory -c web-container --sort PODNAME

# List container memory info from all pods where label app matches web
kubectl-ice memory -l app=web

# List container memory info from all pods where the pod label app is either web or mail
kubectl-ice memory -l "app in (web,mail)"
```
## Ports

``` shell
Usage:
  kubectl-ice ports [flags]

Aliases:
  ports, port, po

Flags:
      --node-label string              Show the selected node label as a column
  -o, --output string                  Output format, currently csv, list, json and yaml are supported
      --pod-label string               Show the selected pod label as a column
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
      --show-namespace                 Show the namespace column
      --show-node                      Show the node name column
  -T, --show-type                      Show the container type column
  -t, --tree                           Display tree like view instead of the standard list


More information at: https://www.github.com/NimbleArchitect/kubectl-ice

```
also includes standard common kubectl flags

### Examples
``` shell
# List containers port info from pods
kubectl-ice ports

# List container port info from pods output in JSON format
kubectl-ice ports -o json

# List container port info from a single pod
kubectl-ice ports my-pod-4jh36

# List port info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice ports -c web-container

# List port info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice ports -c web-container --sort '!CONTAINER'

# List port info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice ports -c web-container --sort PODNAME

# List container port info from all pods where label app matches web
kubectl-ice ports -l app=web

# List container port info from all pods where the pod label app is either web or mail
kubectl-ice ports -l "app in (web,mail)"
```
## Probes

``` shell
Usage:
  kubectl-ice probes [flags]

Aliases:
  probes, probe

Flags:
      --node-label string              Show the selected node label as a column
  -o, --output string                  Output format, currently csv, list, json and yaml are supported
      --pod-label string               Show the selected pod label as a column
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
      --show-namespace                 Show the namespace column
      --show-node                      Show the node name column
  -T, --show-type                      Show the container type column
  -t, --tree                           Display tree like view instead of the standard list


More information at: https://www.github.com/NimbleArchitect/kubectl-ice

```
also includes standard common kubectl flags

### Examples
``` shell
# List containers probe info from pods
kubectl-ice probes

# List container probe info from pods output in JSON format
kubectl-ice probes -o json

# List container probe info from a single pod
kubectl-ice probes my-pod-4jh36

# List probe info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice probes -c web-container

# List probe info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice probes -c web-container --sort '!CONTAINER'

# List probe info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice probes -c web-container --sort PODNAME

# List container probe info from all pods where label app matches web
kubectl-ice probes -l app=web

# List container probe info from all pods where the pod label app is either web or mail
kubectl-ice probes -l "app in (web,mail)"
```
## Restarts

``` shell
Usage:
  kubectl-ice restarts [flags]

Aliases:
  restarts, restart

Flags:
      --node-label string              Show the selected node label as a column
      --oddities                       show only the outlier rows that dont fall within the computed range
  -o, --output string                  Output format, currently csv, list, json and yaml are supported
      --pod-label string               Show the selected pod label as a column
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
      --show-namespace                 Show the namespace column
      --show-node                      Show the node name column
  -T, --show-type                      Show the container type column
  -t, --tree                           Display tree like view instead of the standard list


More information at: https://www.github.com/NimbleArchitect/kubectl-ice

```
also includes standard common kubectl flags

### Examples
``` shell
# List individual container restart count from pods
kubectl-ice restarts

# List conttainers restart count from pods output in JSON format
kubectl-ice restarts -o json

# List restart count from all containers in a single pod
kubectl-ice restarts my-pod-4jh36

# List restart count of all containers named web-container searching all
# pods in the current namespace
kubectl-ice restarts -c web-container

# List restart count of containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice restarts -c web-container --sort '!CONTAINER'

# List restart count of containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice restarts -c web-container --sort PODNAME

# List container restart count from all pods where label app equals web
kubectl-ice restarts -l app=web

# List restart count from all containers where the pod label app is either web or mail
kubectl-ice restarts -l "app in (web,mail)"
```
## Security

``` shell
Usage:
  kubectl-ice security [flags]

Aliases:
  security, sec

Flags:
      --node-label string              Show the selected node label as a column
  -o, --output string                  Output format, currently csv, list, json and yaml are supported
      --pod-label string               Show the selected pod label as a column
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
      --selinux                        show the SELinux context thats applied to the containers
      --show-namespace                 Show the namespace column
      --show-node                      Show the node name column
  -T, --show-type                      Show the container type column
  -t, --tree                           Display tree like view instead of the standard list


More information at: https://www.github.com/NimbleArchitect/kubectl-ice

```
also includes standard common kubectl flags

### Examples
``` shell
# List container security info from pods
kubectl-ice security

# List container security info from pods output in JSON format
kubectl-ice security -o json

# List container security info from a single pod
kubectl-ice security my-pod-4jh36

# List security info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice security -c web-container

# List security info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice security -c web-container --sort '!CONTAINER'

# List security info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice security -c web-container --sort PODNAME

# List container security info from all pods where label app matches web
kubectl-ice security -l app=web

# List container security info from all pods where the pod label app is either web or mail
kubectl-ice security -l "app in (web,mail)"
```
## Status

``` shell
Usage:
  kubectl-ice status [flags]

Aliases:
  status, st

Flags:
  -d, --details                        Display the timestamp instead of age along with the message column
      --node-label string              Show the selected node label as a column
      --oddities                       show only the outlier rows that dont fall within the computed range
  -o, --output string                  Output format, currently csv, list, json and yaml are supported
      --pod-label string               Show the selected pod label as a column
  -p, --previous                       Show previous state
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
      --show-namespace                 Show the namespace column
      --show-node                      Show the node name column
  -T, --show-type                      Show the container type column
  -t, --tree                           Display tree like view instead of the standard list


More information at: https://www.github.com/NimbleArchitect/kubectl-ice

```
also includes standard common kubectl flags

### Examples
``` shell
# List individual container status from pods
kubectl-ice status

# List conttainers status from pods output in JSON format
kubectl-ice status -o json

# List status from all container in a single pod
kubectl-ice status my-pod-4jh36

# List previous container status from a single pod
kubectl-ice status -p my-pod-4jh36

# List status of all containers named web-container searching all
# pods in the current namespace
kubectl-ice status -c web-container

# List status of containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice status -c web-container --sort '!CONTAINER'

# List status of containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice status -c web-container --sort PODNAME

# List container status from all pods where label app equals web
kubectl-ice status -l app=web

# List status from all containers where the pods label app is either web or mail
kubectl-ice status -l "app in (web,mail)"
```
## Volumes

``` shell
Usage:
  kubectl-ice volumes [flags]

Aliases:
  volumes, volume, vol

Flags:
  -d, --device                         show raw block device mappings within a container
      --node-label string              Show the selected node label as a column
  -o, --output string                  Output format, currently csv, list, json and yaml are supported
      --pod-label string               Show the selected pod label as a column
      --select string                  Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != 
      --show-namespace                 Show the namespace column
      --show-node                      Show the node name column
  -T, --show-type                      Show the container type column
  -t, --tree                           Display tree like view instead of the standard list


More information at: https://www.github.com/NimbleArchitect/kubectl-ice

```
also includes standard common kubectl flags

### Examples
``` shell
# List volumes from containers inside pods from current namespace
kubectl-ice volumes

# List volumes from conttainers output in JSON format
kubectl-ice volumes -o json

# List all container volumes from a single pod
kubectl-ice volumes my-pod-4jh36

# List volumes from all containers named web-container searching all
# pods in the current namespace
kubectl-ice volumes -c web-container

# List volumes from container web-container searching all pods in current
# namespace sorted by volume name in descending order (notice the ! charator)
kubectl-ice volumes -c web-container --sort '!VOLUME'

# List volumes from container web-container searching all pods in current
# namespace sorted by volume name in ascending order
kubectl-ice volumes -c web-container --sort MOUNT-POINT

# List container volume info from all pods where label app equals web
kubectl-ice volumes -l app=web

# List volumes from all containers where the pod label app is web or mail
kubectl-ice volumes -l "app in (web,mail)"
```

