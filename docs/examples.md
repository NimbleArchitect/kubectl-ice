### Single pod info
Shows the currently used memory along with the configured memory requests and limits of all containers (side cars) in the pod named web-pod
``` shell
$ kubectl-ice memory web-pod
CONTAINER    USED  REQUEST  LIMIT  %REQ  %LIMIT
app-watcher  0     1M       512M   -     -
app-broken   0     1M       512M   -     -
myapp        0     1M       256M   -     -

```
### Using labels
using labels you can search all pods that are part of a deployment where the label app matches demoprobe and list selected information about the containers in each pod, this example shows the currently configured probe information and gives details of configured startup, readiness and liveness probes of each container
``` shell
$ kubectl-ice probes -l app=demoprobe
PODNAME                      CONTAINER     PROBE     DELAY  PERIOD  TIMEOUT  SUCCESS  FAILURE  CHECK    ACTION
demo-probe-76b66d5766-jlnnd  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-jlnnd  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-jlnnd  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/
demo-probe-76b66d5766-jmqpf  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-jmqpf  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-jmqpf  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/

```
### Named containers
the optional container flag (-c) searchs all selected pods and lists only containers that match the name web-frontend
``` shell
$ kubectl-ice command -c web-frontend
T  PODNAME                           CONTAINER     COMMAND                                      ARGUMENTS
S  demo-memory-7ddb58cd5b-d8f6q      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-dh9zq      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-dq9lq      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-qphl4      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-odd-cpu-5f947f9db4-5w88x     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-6k2wf     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-86d7q     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-9gmhq     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-bsfzf     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-cf787     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-g5k7q     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-h5tql     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-h9zpr     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-kddvb     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-n5slb     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-nmxrj     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-qkvzq     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-tvqbs     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-x24r2     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-xv67t     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-probe-76b66d5766-jlnnd       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
S  demo-probe-76b66d5766-jmqpf       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
S  demo-random-cpu-669b7888b9-8gwn7  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-9rwf9  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-kvbhk  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-wgr46  web-frontend  python /myapp/randomcpuapp.py                -

```
### Labels and containers
you can also search specific pods and list all containers with a specific name, in this example all pods with the label app=userandomcpu are searched and only the containers that match the name web-fronteend are shown
``` shell
$ kubectl-ice cpu -l app=userandomcpu -c web-frontend
PODNAME                           CONTAINER     USED  REQUEST  LIMIT  %REQ  %LIMIT
demo-random-cpu-669b7888b9-8gwn7  web-frontend  0m    125m     1000m  -     -
demo-random-cpu-669b7888b9-9rwf9  web-frontend  0m    125m     1000m  -     -
demo-random-cpu-669b7888b9-kvbhk  web-frontend  0m    125m     1000m  -     -
demo-random-cpu-669b7888b9-wgr46  web-frontend  0m    125m     1000m  -     -

```
### Container status
most commands work the same way including the status command which also lets you see which container(s) are causing the restarts and by using the optional --previous flag you can view the containers previous exit code
``` shell
$ kubectl-ice status -l app=myapp --previous
T  PODNAME  CONTAINER    STATE       REASON              EXIT-CODE  SIGNAL  TIMESTAMP                      MESSAGE
S  web-pod  app-broken   Terminated  Error               1          0       2022-06-06 10:08:54 +0100 BST  -
S  web-pod  app-watcher  Terminated  Error               2          0       2022-06-06 10:09:55 +0100 BST  -
S  web-pod  myapp        Terminated  ContainerCannotRun  127        0       2022-06-06 10:08:26 +0100 BST  OCI runtime create failed: container_linux.go:380: starting container process caused: exec: "python /myapp/mainapp.py\n": stat python /myapp/mainapp.py\n: no such file or directory: unknown
I  web-pod  app-init     -           -                   -          -       -                              -

```
### Container images
need to chack on the currently configured image versions use the image command
``` shell
$ kubectl-ice image -l app=userandomcpu
T  PODNAME                           CONTAINER       PULL          IMAGE
S  demo-random-cpu-669b7888b9-8gwn7  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-8gwn7  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-8gwn7  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-9rwf9  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-9rwf9  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-9rwf9  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-kvbhk  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-kvbhk  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-kvbhk  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-wgr46  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-wgr46  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-wgr46  init-myservice  IfNotPresent  busybox:1.28

```
### Advanced labels
return memory requests size and limits of each container where the pods have an app label that matches useoddcpu and the container name is equal to nginx
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-5w88x  web-frontend  103m  1m       1000m  10285.12  10.29
demo-odd-cpu-5f947f9db4-6k2wf  web-frontend  3m    1m       1000m  216.82    0.22
demo-odd-cpu-5f947f9db4-86d7q  web-frontend  2m    1m       1000m  182.98    0.18
demo-odd-cpu-5f947f9db4-9gmhq  web-frontend  3m    1m       1000m  271.87    0.27
demo-odd-cpu-5f947f9db4-bsfzf  web-frontend  3m    1m       1000m  228.82    0.23
demo-odd-cpu-5f947f9db4-cf787  web-frontend  3m    1m       1000m  211.69    0.21
demo-odd-cpu-5f947f9db4-g5k7q  web-frontend  2m    1m       1000m  180.05    0.18
demo-odd-cpu-5f947f9db4-h5tql  web-frontend  3m    1m       1000m  218.84    0.22
demo-odd-cpu-5f947f9db4-h9zpr  web-frontend  2m    1m       1000m  177.78    0.18
demo-odd-cpu-5f947f9db4-kddvb  web-frontend  2m    1m       1000m  188.46    0.19
demo-odd-cpu-5f947f9db4-n5slb  web-frontend  114m  1m       1000m  11377.78  11.38
demo-odd-cpu-5f947f9db4-nmxrj  web-frontend  2m    1m       1000m  176.04    0.18
demo-odd-cpu-5f947f9db4-qkvzq  web-frontend  2m    1m       1000m  168.35    0.17
demo-odd-cpu-5f947f9db4-tvqbs  web-frontend  2m    1m       1000m  183.71    0.18
demo-odd-cpu-5f947f9db4-x24r2  web-frontend  3m    1m       1000m  228.32    0.23
demo-odd-cpu-5f947f9db4-xv67t  web-frontend  2m    1m       1000m  168.46    0.17

```
### Odditites and sorting
given the listed output above the optional --oddities flag picks out the containers that have a high cpu usage when compared to the other containers listed we also sort the list in descending order by the %REQ column
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend --oddities --sort '!%REQ'
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-n5slb  web-frontend  114m  1m       1000m  11377.78  11.38
demo-odd-cpu-5f947f9db4-5w88x  web-frontend  103m  1m       1000m  10285.12  10.29

```
### Pod volumes
list all container volumes with mount points
``` shell
$ kubectl-ice volumes web-pod
CONTAINER    VOLUME                 TYPE       BACKING           SIZE  RO    MOUNT-POINT
app-init     kube-api-access-tzk9t  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-watcher  kube-api-access-tzk9t  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-broken   kube-api-access-tzk9t  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
myapp        app                    ConfigMap  app.py            -     false /myapp/
myapp        kube-api-access-tzk9t  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount

```
### Pod exec command
retrieves the command line and any arguments specified at the container level
``` shell
$ kubectl-ice command web-pod
T  CONTAINER    COMMAND                     ARGUMENTS
S  app-watcher  python /myapp/mainapp.py    -
S  app-broken   sh -c sleep 2; exit 1       -
S  myapp        python /myapp/mainapp.py\n  -
I  app-init     sh -c sleep 2; exit 0       -

```
### Excluding rows
use the --match flag to show only the output rows where the used memory column is greater than or equal to 1MB, this has the effect of exclusing any row where the used memory column is currently under 1MB, the value 1 can be replace with any whole number in megabytes, to show only used memory greater than 1GB you would replace 1 with 1000
``` shell
$ kubectl-ice mem -l app=userandomcpu --match 'used>=1'
PODNAME                           CONTAINER     USED    REQUEST  LIMIT  %REQ    %LIMIT
demo-random-cpu-669b7888b9-8gwn7  nginx         4.10Mi  1M       256M   429.67  1.68
demo-random-cpu-669b7888b9-9rwf9  nginx         4.00Mi  1M       256M   419.43  1.64
demo-random-cpu-669b7888b9-kvbhk  nginx         4.37Mi  1M       256M   457.93  1.79
demo-random-cpu-669b7888b9-wgr46  nginx         4.24Mi  1M       256M   444.83  1.74

```
### Extra selections
using the --select flag allows you to filter the pod selection to only pods that have a priorityClassName thats equal to system-cluster-critical, you can also match against priority
``` shell
$ kubectl-ice status --select 'priorityClassName=system-cluster-critical' -A
T  PODNAME                          CONTAINER       READY  STARTED  RESTARTS  STATE    REASON  EXIT-CODE  SIGNAL  TIMESTAMP                      MESSAGE
S  coredns-78fcd69978-gtg8c         coredns         true   true     19        Running  -       -          -       2022-06-06 09:12:04 +0100 BST  -
S  metrics-server-77c99ccb96-z86xc  metrics-server  true   true     25        Running  -       -          -       2022-06-06 09:12:04 +0100 BST  -

```
### Security information
listing runAsUser and runAsGroup settings along with other related container security information
``` shell
$ kubectl-ice security -n kube-system
T  PODNAME                           CONTAINER                ALLOW_PRIVILEGE_ESCALATION  PRIVILEGED  RO_ROOT_FS  RUN_AS_NON_ROOT  RUN_AS_USER  RUN_AS_GROUP
S  coredns-78fcd69978-gtg8c          coredns                  false                       -           true        -                -            -
S  etcd-minikube                     etcd                     -                           -           -           -                -            -
S  kube-apiserver-minikube           kube-apiserver           -                           -           -           -                -            -
S  kube-controller-manager-minikube  kube-controller-manager  -                           -           -           -                -            -
S  kube-proxy-4p6q8                  kube-proxy               -                           true        -           -                -            -
S  kube-scheduler-minikube           kube-scheduler           -                           -           -           -                -            -
S  metrics-server-77c99ccb96-z86xc   metrics-server           -                           -           true        true             1000         -
S  storage-provisioner               storage-provisioner      -                           -           -           -                -            -

```
### POSIX capabilities
display configured capabilities related to each container
``` shell
$ kubectl-ice capabilities -n kube-system
T  PODNAME                           CONTAINER                ADD               DROP
S  coredns-78fcd69978-gtg8c          coredns                  NET_BIND_SERVICE  all
S  etcd-minikube                     etcd                     -                 -
S  kube-apiserver-minikube           kube-apiserver           -                 -
S  kube-controller-manager-minikube  kube-controller-manager  -                 -
S  kube-proxy-4p6q8                  kube-proxy               -                 -
S  kube-scheduler-minikube           kube-scheduler           -                 -
S  metrics-server-77c99ccb96-z86xc   metrics-server           -                 -
S  storage-provisioner               storage-provisioner      -                 -

```
