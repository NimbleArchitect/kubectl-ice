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
demo-probe-76b66d5766-g84l9  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-g84l9  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-g84l9  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/
demo-probe-76b66d5766-zgvlh  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-zgvlh  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-zgvlh  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/

```
### Named containers
the optional container flag (-c) searchs all selected pods and lists only containers that match the name web-frontend
``` shell
$ kubectl-ice command -c web-frontend
T  PODNAME                           CONTAINER     COMMAND                                      ARGUMENTS
S  demo-memory-7ddb58cd5b-95xs2      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-qct82      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-st45s      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-sx5kn      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-odd-cpu-5f947f9db4-4dg5q     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-4g2jz     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-5srw5     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-68bps     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-8l8hf     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-9mb5r     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-dd97w     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-gg2nm     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-ggwh5     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-gnvhb     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-hss2g     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-ph6ml     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-pjcm5     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-rshnq     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-wglts     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-probe-76b66d5766-g84l9       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
S  demo-probe-76b66d5766-zgvlh       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
S  demo-random-cpu-669b7888b9-5xmxb  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-b4d48  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-hz62z  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-rmm6s  web-frontend  python /myapp/randomcpuapp.py                -

```
### Labels and containers
you can also search specific pods and list all containers with a specific name, in this example all pods with the label app=userandomcpu are searched and only the containers that match the name web-fronteend are shown
``` shell
$ kubectl-ice cpu -l app=userandomcpu -c web-frontend
PODNAME                           CONTAINER     USED  REQUEST  LIMIT  %REQ    %LIMIT
demo-random-cpu-669b7888b9-5xmxb  web-frontend  0m    125m     1000m  -       -
demo-random-cpu-669b7888b9-b4d48  web-frontend  0m    125m     1000m  -       -
demo-random-cpu-669b7888b9-hz62z  web-frontend  542m  125m     1000m  433.09  54.14
demo-random-cpu-669b7888b9-rmm6s  web-frontend  0m    125m     1000m  -       -

```
### Container status
most commands work the same way including the status command which also lets you see which container(s) are causing the restarts and by using the optional --previous flag you can view the containers previous exit code
``` shell
$ kubectl-ice status -l app=myapp --previous
T  PODNAME  CONTAINER    STATE       REASON              EXIT-CODE  SIGNAL  TIMESTAMP                      MESSAGE
S  web-pod  app-broken   Terminated  Error               1          0       2022-06-14 17:06:49 +0100 BST  -
S  web-pod  app-watcher  Terminated  Error               2          0       2022-06-14 17:07:07 +0100 BST  -
S  web-pod  myapp        Terminated  ContainerCannotRun  127        0       2022-06-14 17:05:37 +0100 BST  OCI runtime create failed: container_linux.go:380: starting container process caused: exec: "python /myapp/mainapp.py\n": stat python /myapp/mainapp.py\n: no such file or directory: unknown
I  web-pod  app-init     -           -                   -          -       -                              -

```
### Container images
need to chack on the currently configured image versions use the image command
``` shell
$ kubectl-ice image -l app=userandomcpu
T  PODNAME                           CONTAINER       PULL          IMAGE
S  demo-random-cpu-669b7888b9-5xmxb  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-5xmxb  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-5xmxb  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-b4d48  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-b4d48  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-b4d48  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-hz62z  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-hz62z  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-hz62z  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-rmm6s  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-rmm6s  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-rmm6s  init-myservice  IfNotPresent  busybox:1.28

```
### Advanced labels
return memory requests size and limits of each container where the pods have an app label that matches useoddcpu and the container name is equal to nginx
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-4dg5q  web-frontend  3m    1m       1000m  228.80    0.23
demo-odd-cpu-5f947f9db4-4g2jz  web-frontend  3m    1m       1000m  272.17    0.27
demo-odd-cpu-5f947f9db4-5srw5  web-frontend  135m  1m       1000m  13421.78  13.42
demo-odd-cpu-5f947f9db4-68bps  web-frontend  3m    1m       1000m  218.85    0.22
demo-odd-cpu-5f947f9db4-8l8hf  web-frontend  132m  1m       1000m  13132.87  13.13
demo-odd-cpu-5f947f9db4-9mb5r  web-frontend  3m    1m       1000m  231.19    0.23
demo-odd-cpu-5f947f9db4-dd97w  web-frontend  3m    1m       1000m  230.31    0.23
demo-odd-cpu-5f947f9db4-gg2nm  web-frontend  3m    1m       1000m  206.96    0.21
demo-odd-cpu-5f947f9db4-ggwh5  web-frontend  3m    1m       1000m  235.47    0.24
demo-odd-cpu-5f947f9db4-gnvhb  web-frontend  3m    1m       1000m  263.84    0.26
demo-odd-cpu-5f947f9db4-hss2g  web-frontend  3m    1m       1000m  210.29    0.21
demo-odd-cpu-5f947f9db4-ph6ml  web-frontend  3m    1m       1000m  218.21    0.22
demo-odd-cpu-5f947f9db4-pjcm5  web-frontend  3m    1m       1000m  259.03    0.26
demo-odd-cpu-5f947f9db4-rshnq  web-frontend  3m    1m       1000m  216.83    0.22
demo-odd-cpu-5f947f9db4-wglts  web-frontend  3m    1m       1000m  236.64    0.24

```
### Odditites and sorting
given the listed output above the optional --oddities flag picks out the containers that have a high cpu usage when compared to the other containers listed we also sort the list in descending order by the %REQ column
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend --oddities --sort '!%REQ'
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-5srw5  web-frontend  135m  1m       1000m  13421.78  13.42
demo-odd-cpu-5f947f9db4-8l8hf  web-frontend  132m  1m       1000m  13132.87  13.13

```
### Pod volumes
list all container volumes with mount points
``` shell
$ kubectl-ice volumes web-pod
CONTAINER    VOLUME                 TYPE       BACKING           SIZE  RO    MOUNT-POINT
app-init     kube-api-access-rjpxb  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-watcher  kube-api-access-rjpxb  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-broken   kube-api-access-rjpxb  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
myapp        app                    ConfigMap  app.py            -     false /myapp/
myapp        kube-api-access-rjpxb  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount

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
demo-random-cpu-669b7888b9-5xmxb  nginx         2.54Mi  1M       256M   265.83  1.04
demo-random-cpu-669b7888b9-b4d48  nginx         2.44Mi  1M       256M   255.59  1.00
demo-random-cpu-669b7888b9-hz62z  web-frontend  7.61Mi  1M       256M   797.49  3.12
demo-random-cpu-669b7888b9-hz62z  nginx         3.00Mi  1M       256M   314.16  1.23
demo-random-cpu-669b7888b9-rmm6s  nginx         2.79Mi  1M       256M   292.86  1.14

```
### Extra selections
using the --select flag allows you to filter the pod selection to only pods that have a priorityClassName thats equal to system-cluster-critical, you can also match against priority
``` shell
$ kubectl-ice status --select 'priorityClassName=system-cluster-critical' -A
T  NAMESPACE    PODNAME                          CONTAINER       READY  STARTED  RESTARTS  STATE    REASON  EXIT-CODE  SIGNAL  TIMESTAMP                      MESSAGE
S  kube-system  coredns-78fcd69978-gtg8c         coredns         true   true     23        Running  -       -          -       2022-06-14 09:11:07 +0100 BST  -
S  kube-system  metrics-server-77c99ccb96-z86xc  metrics-server  true   true     30        Running  -       -          -       2022-06-14 09:11:37 +0100 BST  -

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
