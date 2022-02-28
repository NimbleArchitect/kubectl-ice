# kubectl-ice

This plugin shows useful information about the containers inside a pod

# Installation

## From binary
- download the required binary from the release page
- unzip and copy the kubectl-ice file to your path
- run kubectl ice help to check its working

## From Source

```shell
go get https://github.com/NimbleArchitect/kubectl-ice
```

## Usage

The following command are available for `kubectl ice`
```
kubectl ice cpu        # return cpu requests size and limits of each container
kubectl ice help       # Shows the help screen
kubectl ice image      # list the image name and pull status for each container
kubectl ice ip         # list ip addresses of all pods in the namespace listed
kubectl ice memory     # return memory requests size and limits of each container
kubectl ice restarts   # show restart counts for each container in a named pod
kubectl ice stats      # list resource usage of each container in a pod
kubectl ice status     # list status of each container in a pod
kubectl ice volumes    # list all container volumes with mount points
```

### CPU
shows the configured CPU resource requests and limits of each container

``` shell
Usage:
  ice cpu [flags]

Flags:
  -h, --help   help for cpu
```

#### Example
```shell
$ kubectl ice cpu myapp
T  NAME         REQUEST  LIMIT
S  app-watcher  20m      50m
S  app-broken   20m      50m
S  myapp        500m     1
I  app-init     0        0
```

### Image
list the image name and pull status for each container

``` shell
Usage:
  ice image [flags]

Aliases:
  image, im

Flags:
  -h, --help   help for image
```

#### Example
```shell
$ kubectl ice image myapp
T  NAME        PULL          IMAGE
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

Flags:
  -h, --help   help for ip
```

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
  -h, --help   help for memory
```

#### Example
```shell
$ kubectl ice memory myapp
T  NAME         REQUEST  LIMIT
S  app-watcher  500Mi    800Mi
S  app-broken   500Mi    800Mi
S  myapp        500Mi    800Mi
I  app-init     0        0
```

### Restarts
show restart counts for each container in a named pod

``` shell
Usage:
  ice restarts [flags]

Aliases:
  restarts, restart

Flags:
  -h, --help   help for restarts
```

#### Example
```shell
$kubectl ice restarts myapp
T  NAME        RESTARTS
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
  -h, --help   help for stats
  -r, --raw    show raw uncooked values
```

#### Example
```shell
$ kubectl ice stats myapp   
NAME         USED_CPU  CPU_%_REQ  CPU_%_LIMIT  USED_MEM  MEM_%_REQ  MEM_%_LIMIT  
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
  -h, --help       help for status
  -p, --previous   show previous state
```

#### Example
```shell
$ kubectl ice status myapp
T  NAME         READY STARTED  RESTARTS  STATE       REASON     EXIT-CODE  SIGNAL  TIMESTAMP                      MESSAGE  
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

Flags:
  -h, --help   help for volumes
```

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
