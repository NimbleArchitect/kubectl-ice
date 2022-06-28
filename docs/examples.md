### Single pod info
Shows the currently used memory along with the configured memory requests and limits of all containers (side cars) in the pod named web-pod
``` shell
$ kubectl-ice memory web-pod
CONTAINER    USED    REQUEST  LIMIT  %REQ    %LIMIT
app-watcher  5.61Mi  1M       512M   587.78  1.15
app-broken   0       1M       512M   -       -
myapp        5.65Mi  1M       256M   592.28  2.31

```
### Using labels
using labels you can search all pods that are part of a deployment where the label app matches demoprobe and list selected information about the containers in each pod, this example shows the currently configured probe information and gives details of configured startup, readiness and liveness probes of each container
``` shell
$ kubectl-ice probes -l app=demoprobe
PODNAME                      CONTAINER     PROBE     DELAY  PERIOD  TIMEOUT  SUCCESS  FAILURE  CHECK    ACTION
demo-probe-76b66d5766-5qd5c  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-5qd5c  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-5qd5c  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/
demo-probe-76b66d5766-w9jwt  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-w9jwt  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-w9jwt  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/

```
### Named containers
the optional container flag (-c) searchs all selected pods and lists only containers that match the name web-frontend
``` shell
$ kubectl-ice command -c web-frontend
T  PODNAME                           CONTAINER     COMMAND                                      ARGUMENTS
S  demo-memory-7ddb58cd5b-djdsl      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-drl6s      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-jq62c      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-qfjps      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-odd-cpu-5f947f9db4-5cjv6     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-7bmkw     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-bb2d5     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-bp9qm     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-dgxj6     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-f4vf5     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-gb4xf     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-hxdfq     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-j5lrb     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-k25kl     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-prxrf     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-tvrx4     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-v77vj     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-wqgqz     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-xf298     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-zjvxk     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-probe-76b66d5766-5qd5c       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
S  demo-probe-76b66d5766-w9jwt       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
S  demo-random-cpu-669b7888b9-jqrjq  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-mvd4m  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-qmc5b  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-xp8dk  web-frontend  python /myapp/randomcpuapp.py                -

```
### Alternate status view
the tree flag shows the containers and pods in a tree like view
``` shell
$ kubectl-ice status -l app=myapp --tree
NAMESPACE  NAME                      READY  STARTED  RESTARTS  STATE       REASON            EXIT-CODE  SIGNAL  AGE
ice        Pod/web-pod               -      -        0         Running     -                 -          -       52m
ice        └─Container/app-broken    false  false    14        Waiting     CrashLoopBackOff  -          -       -
ice        └─Container/app-watcher   true   true     0         Running     -                 -          -       51m
ice        └─Container/myapp         true   true     0         Running     -                 -          -       51m
ice        └─InitContainer/app-init  true   -        0         Terminated  Completed         0          0       52m

```
### Labels and containers
you can also search specific pods and list all containers with a specific name, in this example all pods with the label app=userandomcpu are searched and only the containers that match the name web-fronteend are shown
``` shell
$ kubectl-ice cpu -l app=userandomcpu -c web-frontend
PODNAME                           CONTAINER     USED  REQUEST  LIMIT  %REQ    %LIMIT
demo-random-cpu-669b7888b9-jqrjq  web-frontend  326m  125m     1000m  260.59  32.57
demo-random-cpu-669b7888b9-mvd4m  web-frontend  392m  125m     1000m  312.83  39.10
demo-random-cpu-669b7888b9-qmc5b  web-frontend  0m    125m     1000m  -       -
demo-random-cpu-669b7888b9-xp8dk  web-frontend  143m  125m     1000m  114.40  14.30

```
### Status details
using the details flag displays the timestamp and message columns
``` shell
$ kubectl-ice status -l app=myapp --details
T  PODNAME  CONTAINER    READY  STARTED  RESTARTS  STATE       REASON            EXIT-CODE  SIGNAL  TIMESTAMP            MESSAGE
S  web-pod  app-broken   false  false    14        Waiting     CrashLoopBackOff  -          -       -                    back-off 5m0s restarting failed
S  web-pod  app-watcher  true   true     0         Running     -                 -          -       2022-06-28 18:44:23  -
S  web-pod  myapp        true   true     0         Running     -                 -          -       2022-06-28 18:44:34  -
I  web-pod  app-init     true   -        0         Terminated  Completed         0          0       2022-06-28 18:43:49  -

```
### Container status
most commands work the same way including the status command which also lets you see which container(s) are causing the restarts and by using the optional --previous flag you can view the containers previous exit code
``` shell
$ kubectl-ice status -l app=myapp --previous
T  PODNAME  CONTAINER    STATE       REASON  EXIT-CODE  SIGNAL  TIMESTAMP            AGE    MESSAGE
S  web-pod  app-broken   Terminated  Error   1          0       2022-06-28 19:33:52  2m24s  -
S  web-pod  app-watcher  -           -       -          -       -                    292y   -
S  web-pod  myapp        -           -       -          -       -                    292y   -
I  web-pod  app-init     -           -       -          -       -                    292y   -

```
### Container images
need to chack on the currently configured image versions use the image command
``` shell
$ kubectl-ice image -l app=userandomcpu
T  PODNAME                           CONTAINER       PULL          IMAGE
S  demo-random-cpu-669b7888b9-jqrjq  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-jqrjq  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-jqrjq  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-mvd4m  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-mvd4m  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-mvd4m  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-qmc5b  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-qmc5b  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-qmc5b  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-xp8dk  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-xp8dk  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-xp8dk  init-myservice  IfNotPresent  busybox:1.28

```
### Advanced labels
return memory requests size and limits of each container where the pods have an app label that matches useoddcpu and the container name is equal to nginx
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-5cjv6  web-frontend  135m  1m       1000m  13433.33  13.43
demo-odd-cpu-5f947f9db4-7bmkw  web-frontend  3m    1m       1000m  209.55    0.21
demo-odd-cpu-5f947f9db4-bb2d5  web-frontend  3m    1m       1000m  267.87    0.27
demo-odd-cpu-5f947f9db4-bp9qm  web-frontend  3m    1m       1000m  220.67    0.22
demo-odd-cpu-5f947f9db4-dgxj6  web-frontend  3m    1m       1000m  212.91    0.21
demo-odd-cpu-5f947f9db4-f4vf5  web-frontend  135m  1m       1000m  13449.89  13.45
demo-odd-cpu-5f947f9db4-gb4xf  web-frontend  3m    1m       1000m  272.25    0.27
demo-odd-cpu-5f947f9db4-hxdfq  web-frontend  2m    1m       1000m  190.80    0.19
demo-odd-cpu-5f947f9db4-j5lrb  web-frontend  126m  1m       1000m  12534.64  12.53
demo-odd-cpu-5f947f9db4-k25kl  web-frontend  3m    1m       1000m  241.68    0.24
demo-odd-cpu-5f947f9db4-prxrf  web-frontend  2m    1m       1000m  187.68    0.19
demo-odd-cpu-5f947f9db4-tvrx4  web-frontend  2m    1m       1000m  188.06    0.19
demo-odd-cpu-5f947f9db4-v77vj  web-frontend  3m    1m       1000m  262.66    0.26
demo-odd-cpu-5f947f9db4-wqgqz  web-frontend  3m    1m       1000m  229.98    0.23
demo-odd-cpu-5f947f9db4-xf298  web-frontend  2m    1m       1000m  179.28    0.18
demo-odd-cpu-5f947f9db4-zjvxk  web-frontend  3m    1m       1000m  219.51    0.22

```
### Odditites and sorting
given the listed output above the optional --oddities flag picks out the containers that have a high cpu usage when compared to the other containers listed we also sort the list in descending order by the %REQ column
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend --oddities --sort '!%REQ'
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-f4vf5  web-frontend  135m  1m       1000m  13449.89  13.45
demo-odd-cpu-5f947f9db4-5cjv6  web-frontend  135m  1m       1000m  13433.33  13.43
demo-odd-cpu-5f947f9db4-j5lrb  web-frontend  126m  1m       1000m  12534.64  12.53

```
### Pod volumes
list all container volumes with mount points
``` shell
$ kubectl-ice volumes web-pod
CONTAINER    VOLUME                 TYPE       BACKING           SIZE  RO    MOUNT-POINT
app-init     kube-api-access-g8487  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-watcher  app                    ConfigMap  app.py            -     false /myapp/
app-watcher  kube-api-access-g8487  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-broken   kube-api-access-g8487  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
myapp        app                    ConfigMap  app.py            -     false /myapp/
myapp        kube-api-access-g8487  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount

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
demo-random-cpu-669b7888b9-jqrjq  web-frontend  8.22Mi  1M       256M   862.21  3.37
demo-random-cpu-669b7888b9-jqrjq  nginx         2.52Mi  1M       256M   264.60  1.03
demo-random-cpu-669b7888b9-mvd4m  web-frontend  8.68Mi  1M       256M   909.72  3.55
demo-random-cpu-669b7888b9-mvd4m  nginx         2.46Mi  1M       256M   258.46  1.01
demo-random-cpu-669b7888b9-qmc5b  web-frontend  0.89Mi  1M       256M   93.80   0.37
demo-random-cpu-669b7888b9-qmc5b  nginx         2.55Mi  1M       256M   267.06  1.04
demo-random-cpu-669b7888b9-xp8dk  web-frontend  8.27Mi  1M       256M   866.71  3.39
demo-random-cpu-669b7888b9-xp8dk  nginx         2.55Mi  1M       256M   267.47  1.04

```
### Extra selections
using the --select flag allows you to filter the pod selection to only pods that have a priorityClassName thats equal to system-cluster-critical, you can also match against priority
``` shell
$ kubectl-ice status --select 'priorityClassName=system-cluster-critical' -A
T  NAMESPACE    PODNAME                          CONTAINER       READY  STARTED  RESTARTS  STATE    REASON  EXIT-CODE  SIGNAL  AGE
S  kube-system  coredns-78fcd69978-gtg8c         coredns         true   true     26        Running  -       -          -       8d
S  kube-system  metrics-server-77c99ccb96-z86xc  metrics-server  true   true     33        Running  -       -          -       8d

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
