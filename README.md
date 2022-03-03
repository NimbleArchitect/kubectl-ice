## *Still in development*

# kubectl-ice

This plugin shows useful information about the containers inside a pod useful for trouble shooting container issues


Note: * *requires metrics server* *

# Installation

## From binary
- download the required binary from the release page
- unzip and copy the kubectl-ice file to your path
- run kubectl ice help to check its working

## From Source

```shell
git clone https://github.com/NimbleArchitect/kubectl-ice.git
make bin
```

## Usage

The following command are available for `kubectl ice`
```
kubectl ice command    # retrieves the command line and any arguments specified at the container level
kubectl ice cpu        # return cpu requests size, limits and usage of each container
kubectl ice help       # Shows the help screen
kubectl ice image      # list the image name and pull status for each container
kubectl ice ip         # list ip addresses of all pods in the namespace listed
kubectl ice memory     # return memory requests size, limits and usage of each container
kubectl ice ports      # shows ports exposed by the containers in a pod
kubectl ice probes     # shows details of configured startup, readiness and liveness probes of each container
kubectl ice restarts   # show restart counts for each container in a named pod
kubectl ice status     # list status of each container in a pod
kubectl ice volumes    # list all container volumes with mount points
```


### Command
retrieves the command line and any arguments specified at the container level

``` shell
Usage:
  ice command [flags]

Aliases:
  command, cmd, exec, args
```
also includes standard common kubectl flags

#### Example
```shell
$ kubectl ice command myapp
T  CONTAINER     COMMAND    ARGUMENTS
S  app-watcher   -          -
S  app-broken    /bin/bash  -s exit 1
S  myapp         -          -
I  app-init      init.sh    -
```
### CPU
shows the configured CPU resource requests and limits of each container

``` shell
Usage:
  ice cpu [flags]

Flags:
  -r, --raw              show raw uncooked values
```
also includes standard common kubectl flags

#### Example
```shell
$ kubectl ice cpu myapp
T  CONTAINER    USED  REQUEST  LIMIT   %REQ  %LIMIT
S  app-watcher  0     20m      50m     0     0
S  app-broken   0     20m      50m     0     0
S  myapp        1     500m     1       200   100
I  app-init     0        0     0       0     0
```

### Image
list the image name and pull status for each container

``` shell
Usage:
  ice image [flags]

Aliases:
  image, im
```
also includes standard common kubectl flags

#### Example
```shell
$ kubectl ice image myapp
T  CONTAINER   PULL          IMAGE
S  app-watcher Always        amouat/network-utils
S  app-broken  IfNotPresent  busybox:1.28
S  myapp       Always        amouat/network-utils
I  app-init    Always        amouat/network-utils
```

### IP
list ip addresses of all pods in the namespace listed

``` shell
Usage:
  ice ip [flags]
```
also includes standard common kubectl flags

#### Example
```shell
$ kubectl ice ip myapp   
NAME  IP
myapp 172.17.0.2
```

### Memory
return memory requests size and limits of each container

``` shell
Usage:
  ice memory [flags]

Aliases:
  memory, mem

Flags:
  -r, --raw              show raw uncooked values
```
also includes standard common kubectl flags

#### Example
```shell
$ kubectl ice memory myapp
T  CONTAINER    USED  REQUEST  LIMIT   %REQ  %LIMIT
S  app-watcher  0     500Mi    800Mi   0     0
S  app-broken   0     500Mi    800Mi   0     0
S  myapp        1     500Mi    800Mi   0.12  0
I  app-init     0        0     -       -     -
```

### Ports
shows ports exposed by the containers in a pod

``` shell
Usage:
  ice ports [flags]

Aliases:
  ports, port, po
```
also includes standard common kubectl flags

#### Example
```shell
$ kubectl ice ports myapp
T  CONTAINER    PORTNAME  PORT  PROTO  HOSTPORT 
S  app-broken   -         8000  TCP    
S  app-watcher  -         8080  TCP    
S  myapp        http      8080  TCP    
S  keycloak     https     8443  TCP
```

### Probes
shows details of configured startup, readiness and liveness probes of each container

Usage:
  ice probes [flags]

Aliases:
  probes, probe
```
also includes standard common kubectl flags

#### Example
```shell
$ kubectl ice probes myapp
CONTAINER     PROBE     DELAY  PERIOD  TIMEOUT  SUCCESS  FAILURE  CHECK    ACTION
myapp         liveness  0      10      1        1        3        HTTPGet  http://:http/health
app-broken    liveness  0      10      1        1        3        HTTPGet  http://:http/health
```

### Restarts
show restart counts for each container in a named pod

``` shell
Usage:
  ice restarts [flags]

Aliases:
  restarts, restart
```
also includes standard common kubectl flags

#### Example
```shell
$ kubectl ice restarts myapp
T  CONTAINER   RESTARTS
S  app-broken  0
S  app-watcher 0
S  myapp       0
I  app-init    0
```

### Stats
list resource usage of each container in a pod

``` shell
Usage:
  ice stats [flags]

Aliases:
  stats, top, ps

Flags:
  -r, --raw              show raw uncooked values
```
also includes standard common kubectl flags

#### Example
```shell
$ kubectl ice stats myapp   
CONTAINER    USED_CPU  CPU_%_REQ  CPU_%_LIMIT  USED_MEM  MEM_%_REQ  MEM_%_LIMIT  
app-init     0         0          0            0         0          0
app-watcher  0         0.00       0.00         0.92Mi    0.18       0.12
app-broken   0         0.00       0.00         3.95Mi    0.79       0.49
myapp        34        6.673187   3.336594     0.88Mi    0.18       0.11
```

### Status
list current running status of each container in a pod

``` shell
Usage:
  ice status [flags]

Aliases:
  status, st

Flags:
  -p, --previous         show previous state
```
also includes standard common kubectl flags

```

#### Example
```shell
$ kubectl ice status myapp
T  CONTAINER    READY STARTED  RESTARTS  STATE       REASON     EXIT-CODE  SIGNAL  TIMESTAMP                      MESSAGE  
S  app-broken   true  true     0         Running                                   2022-02-28 11:04:24 +0000 GMT           
S  app-watcher  true  true     0         Running                                   2022-02-28 11:04:24 +0000 GMT           
S  myapp        true  true     0         Running                                   2022-02-28 11:04:26 +0000 GMT           
I  app-init     true           0         Terminated  Completed  0          0       2022-02-28 11:04:17 +0000 GMT           

```

### Volumes
list all container volumes with mount points

``` shell
Usage:
  ice volumes [flags]

Aliases:
  volumes, volume, vol
```
also includes standard common kubectl flags


#### Example
```shell
$ kubectl ice volumes myapp
CONTAINER    VOLUME                 TYPE      BACKING SIZE  RO     MOUNT-POINT                                    
app-init     kube-api-access-k7hvs  Projected               true   /var/run/secrets/kubernetes.io/serviceaccount  
app-watcher  appsafe                EmptyDir  Memory        false  /mnt/appsafe                                   
app-watcher  work                   EmptyDir  Memory        false  /mnt/work                                      
app-watcher  shareme                EmptyDir  Memory        false  /etc/shareme                                   
app-watcher  kube-api-access-k7hvs  Projected               true   /var/run/secrets/kubernetes.io/serviceaccount  
app-broken   work                   EmptyDir  Memory        false  /mnt/work                                      
app-broken   appsafe                EmptyDir  Memory        true   /mnt/appsafe                                   
app-broken   kube-api-access-k7hvs  Projected               true   /var/run/secrets/kubernetes.io/serviceaccount  
myapp        appsafe                EmptyDir  Memory        true   /mnt/appsafe                                   
myapp        work                   EmptyDir  Memory        true   /mnt/work                                      
myapp        kube-api-access-k7hvs  Projected               true   /var/run/secrets/kubernetes.io/serviceaccount  
```

## Contributing

All feedback and contributions are welcome, if you want to raise an issue or help with fixes or features please [raise an issue to discuss](https://github.com/NimbleArchitect/kubectl-ice/issues)


## License
Licensed under Apache 2.0 see [LICENSE](https://github.com/NimbleArchitect/kubectl-pod/blob/main/LICENSE)
