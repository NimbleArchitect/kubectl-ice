### Single pod info
Shows the currently used memory along with the configured memory requests and limits of all containers (side cars) in the pod named web-pod
``` shell
$ kubectl-ice memory web-pod
CONTAINER    USED    REQUEST  LIMIT  %REQ    %LIMIT
app-watcher  0       1M       512M   -       -
app-broken   0       1M       512M   -       -
myapp        8.16Mi  1M       256M   855.24  3.34

```
### Using labels
using labels you can search all pods that are part of a deployment where the label app matches demoprobe and list selected information about the containers in each pod, this example shows the currently configured probe information and gives details of configured startup, readiness and liveness probes of each container
``` shell
$ kubectl-ice probes -l app=demoprobe
PODNAME                      CONTAINER     PROBE     DELAY  PERIOD  TIMEOUT  SUCCESS  FAILURE  CHECK    ACTION
demo-probe-76b66d5766-gb6bp  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-gb6bp  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-gb6bp  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/
demo-probe-76b66d5766-jqq99  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-jqq99  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-jqq99  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/

```
### Named containers
the optional container flag (-c) searchs all selected pods and lists only containers that match the name web-frontend
``` shell
$ kubectl-ice command -c web-frontend
T  PODNAME                           CONTAINER     COMMAND                                      ARGUMENTS
S  demo-memory-7ddb58cd5b-4lm9n      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-7fgbf      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-cs757      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-qldjq      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-odd-cpu-5f947f9db4-258br     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-542ql     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-68qxt     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-czkzk     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-f77vp     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-frt4z     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-g84q7     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-mkmvb     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-p6jk5     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-r2rzr     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-t2bj8     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-t8nrf     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-v2cvs     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-vhgk5     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-xshls     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-zvp5z     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-probe-76b66d5766-gb6bp       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
S  demo-probe-76b66d5766-jqq99       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
S  demo-random-cpu-669b7888b9-6jj22  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-lf6jm  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-vq66q  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-z9zwk  web-frontend  python /myapp/randomcpuapp.py                -

```
### Alternate status view
the tree flag show the containers and pods in a tree like view
``` shell
$ kubectl-ice status -l app=myapp --tree
NAMESPACE  NAME                      READY  STARTED  RESTARTS  STATE       REASON            AGE
ice        Pod/web-pod               -      -        0         Running     -                 10m
ice        └─Container/app-broken    false  false    6         Waiting     CrashLoopBackOff  -
ice        └─Container/app-watcher   false  false    6         Waiting     CrashLoopBackOff  -
ice        └─Container/myapp         true   true     0         Running     -                 10m
ice        └─InitContainer/app-init  true   -        0         Terminated  Completed         10m

```
### Labels and containers
you can also search specific pods and list all containers with a specific name, in this example all pods with the label app=userandomcpu are searched and only the containers that match the name web-fronteend are shown
``` shell
$ kubectl-ice cpu -l app=userandomcpu -c web-frontend
PODNAME                           CONTAINER     USED  REQUEST  LIMIT  %REQ  %LIMIT
demo-random-cpu-669b7888b9-6jj22  web-frontend  0m    125m     1000m  -     -
demo-random-cpu-669b7888b9-lf6jm  web-frontend  0m    125m     1000m  -     -
demo-random-cpu-669b7888b9-vq66q  web-frontend  0m    125m     1000m  -     -
demo-random-cpu-669b7888b9-z9zwk  web-frontend  0m    125m     1000m  -     -

```
### Status details
using the details flag displays the timestamp and message columns
``` shell
$ kubectl-ice status -l app=myapp --details
T  PODNAME  CONTAINER    READY  STARTED  RESTARTS  STATE       REASON            EXIT-CODE  SIGNAL  TIMESTAMP            MESSAGE
S  web-pod  app-broken   false  false    6         Waiting     CrashLoopBackOff  -          -       -                    back-off 5m0s restarting failed
S  web-pod  app-watcher  false  false    6         Waiting     CrashLoopBackOff  -          -       -                    back-off 5m0s restarting failed
S  web-pod  myapp        true   true     0         Running     -                 -          -       2022-06-22 22:52:05  -
I  web-pod  app-init     true   -        0         Terminated  Completed         0          0       2022-06-22 22:51:49  -

```
### Container status
most commands work the same way including the status command which also lets you see which container(s) are causing the restarts and by using the optional --previous flag you can view the containers previous exit code
``` shell
$ kubectl-ice status -l app=myapp --previous
T  PODNAME  CONTAINER    STATE       REASON  EXIT-CODE  SIGNAL  AGE
S  web-pod  app-broken   Terminated  Error   1          0       3m49s
S  web-pod  app-watcher  Terminated  Error   2          0       2m23s
S  web-pod  myapp        -           -       -          -       292y
I  web-pod  app-init     -           -       -          -       292y

```
### Container images
need to chack on the currently configured image versions use the image command
``` shell
$ kubectl-ice image -l app=userandomcpu
T  PODNAME                           CONTAINER       PULL          IMAGE
S  demo-random-cpu-669b7888b9-6jj22  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-6jj22  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-6jj22  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-lf6jm  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-lf6jm  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-lf6jm  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-vq66q  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-vq66q  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-vq66q  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-z9zwk  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-z9zwk  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-z9zwk  init-myservice  IfNotPresent  busybox:1.28

```
### Advanced labels
return memory requests size and limits of each container where the pods have an app label that matches useoddcpu and the container name is equal to nginx
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-258br  web-frontend  2m    1m       1000m  181.14    0.18
demo-odd-cpu-5f947f9db4-542ql  web-frontend  126m  1m       1000m  12549.38  12.55
demo-odd-cpu-5f947f9db4-68qxt  web-frontend  2m    1m       1000m  167.06    0.17
demo-odd-cpu-5f947f9db4-czkzk  web-frontend  3m    1m       1000m  201.11    0.20
demo-odd-cpu-5f947f9db4-f77vp  web-frontend  2m    1m       1000m  189.62    0.19
demo-odd-cpu-5f947f9db4-frt4z  web-frontend  2m    1m       1000m  168.54    0.17
demo-odd-cpu-5f947f9db4-g84q7  web-frontend  3m    1m       1000m  233.97    0.23
demo-odd-cpu-5f947f9db4-mkmvb  web-frontend  3m    1m       1000m  215.99    0.22
demo-odd-cpu-5f947f9db4-p6jk5  web-frontend  2m    1m       1000m  192.40    0.19
demo-odd-cpu-5f947f9db4-r2rzr  web-frontend  129m  1m       1000m  12890.36  12.89
demo-odd-cpu-5f947f9db4-t2bj8  web-frontend  2m    1m       1000m  175.11    0.18
demo-odd-cpu-5f947f9db4-t8nrf  web-frontend  2m    1m       1000m  178.12    0.18
demo-odd-cpu-5f947f9db4-v2cvs  web-frontend  3m    1m       1000m  227.29    0.23
demo-odd-cpu-5f947f9db4-vhgk5  web-frontend  3m    1m       1000m  201.97    0.20
demo-odd-cpu-5f947f9db4-xshls  web-frontend  117m  1m       1000m  11644.62  11.64
demo-odd-cpu-5f947f9db4-zvp5z  web-frontend  2m    1m       1000m  181.60    0.18

```
### Odditites and sorting
given the listed output above the optional --oddities flag picks out the containers that have a high cpu usage when compared to the other containers listed we also sort the list in descending order by the %REQ column
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend --oddities --sort '!%REQ'
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-r2rzr  web-frontend  129m  1m       1000m  12890.36  12.89
demo-odd-cpu-5f947f9db4-542ql  web-frontend  126m  1m       1000m  12549.38  12.55
demo-odd-cpu-5f947f9db4-xshls  web-frontend  117m  1m       1000m  11644.62  11.64

```
### Pod volumes
list all container volumes with mount points
``` shell
$ kubectl-ice volumes web-pod
CONTAINER    VOLUME                 TYPE       BACKING           SIZE  RO    MOUNT-POINT
app-init     kube-api-access-f5kt6  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-watcher  kube-api-access-f5kt6  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-broken   kube-api-access-f5kt6  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
myapp        app                    ConfigMap  app.py            -     false /myapp/
myapp        kube-api-access-f5kt6  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount

```
### Pod exec command
retrieves the command line and any arguments specified at the container level
``` shell
$ kubectl-ice command web-pod
T  CONTAINER    COMMAND                   ARGUMENTS
S  app-watcher  python /myapp/mainapp.py  -
S  app-broken   sh -c sleep 2; exit 1     -
S  myapp        python /myapp/mainapp.py  -
I  app-init     sh -c sleep 2; exit 0     -

```
### Excluding rows
use the --match flag to show only the output rows where the used memory column is greater than or equal to 1MB, this has the effect of exclusing any row where the used memory column is currently under 1MB, the value 1 can be replace with any whole number in megabytes, to show only used memory greater than 1GB you would replace 1 with 1000
``` shell
$ kubectl-ice mem -l app=userandomcpu --match 'used>=1'
PODNAME                           CONTAINER     USED    REQUEST  LIMIT  %REQ    %LIMIT
demo-random-cpu-669b7888b9-6jj22  nginx         3.61Mi  1M       256M   378.47  1.48
demo-random-cpu-669b7888b9-lf6jm  nginx         3.55Mi  1M       256M   372.33  1.45
demo-random-cpu-669b7888b9-vq66q  nginx         3.54Mi  1M       256M   371.51  1.45
demo-random-cpu-669b7888b9-z9zwk  nginx         3.45Mi  1M       256M   361.68  1.41

```
### Extra selections
using the --select flag allows you to filter the pod selection to only pods that have a priorityClassName thats equal to system-cluster-critical, you can also match against priority
``` shell
$ kubectl-ice status --select 'priorityClassName=system-cluster-critical' -A
T  NAMESPACE    PODNAME                          CONTAINER       READY  STARTED  RESTARTS  STATE    REASON  EXIT-CODE  SIGNAL  AGE
S  kube-system  coredns-78fcd69978-gtg8c         coredns         true   true     26        Running  -       -          -       2d14h
S  kube-system  metrics-server-77c99ccb96-z86xc  metrics-server  true   true     33        Running  -       -          -       2d14h

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
