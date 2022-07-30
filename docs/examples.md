### Single pod info
Shows the currently used memory along with the configured memory requests and limits of all containers (side cars) in the pod named web-pod
``` shell
$ kubectl-ice memory web-pod
CONTAINER    USED    REQUEST  LIMIT  %REQ    %LIMIT
app-init     0       0        0      -       -
app-watcher  5.25Mi  1M       512M   550.09  1.07
app-broken   0       1M       512M   -       -
myapp        5.23Mi  1M       256M   548.45  2.14

```
### Using labels
using labels you can search all pods that are part of a deployment where the label app matches demoprobe and list selected information about the containers in each pod, this example shows the currently configured probe information and gives details of configured startup, readiness and liveness probes of each container
``` shell
$ kubectl-ice probes -l app=demoprobe
PODNAME                      CONTAINER     PROBE     DELAY  PERIOD  TIMEOUT  SUCCESS  FAILURE  CHECK    ACTION
demo-probe-765fd4d8f7-7s5d7  web-frontend  liveness  10     5       1        1        3        Exec     /bin/true
demo-probe-765fd4d8f7-7s5d7  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-765fd4d8f7-7s5d7  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/
demo-probe-765fd4d8f7-cqr6m  web-frontend  liveness  10     5       1        1        3        Exec     /bin/true
demo-probe-765fd4d8f7-cqr6m  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-765fd4d8f7-cqr6m  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/

```
### Named containers
the optional container flag (-c) searchs all selected pods and lists only containers that match the name web-frontend
``` shell
$ kubectl-ice command -c web-frontend
PODNAME                           CONTAINER     COMMAND                                      ARGUMENTS
demo-memory-7ddb58cd5b-dzp5x      web-frontend  python /myapp/halfmemapp.py                  -
demo-memory-7ddb58cd5b-hxkbt      web-frontend  python /myapp/halfmemapp.py                  -
demo-memory-7ddb58cd5b-wn7tt      web-frontend  python /myapp/halfmemapp.py                  -
demo-memory-7ddb58cd5b-xq2t4      web-frontend  python /myapp/halfmemapp.py                  -
demo-odd-cpu-5f947f9db4-4clnc     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-5z9w2     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-62xjb     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-68f47     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-7hlxl     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-7r5s2     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-8qpl5     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-c5sv6     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-c7scd     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-d6qz6     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-jfwtd     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-lkfs2     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-q85pt     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-qdfnb     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-tz4hj     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-zsmwb     web-frontend  python /myapp/oddcpuapp.py                   -
demo-probe-765fd4d8f7-7s5d7       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
demo-probe-765fd4d8f7-cqr6m       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
demo-random-cpu-55954b64b4-2fnvf  web-frontend  python /myapp/randomcpuapp.py                -
demo-random-cpu-55954b64b4-j5fgx  web-frontend  python /myapp/randomcpuapp.py                -
demo-random-cpu-55954b64b4-kfc5q  web-frontend  python /myapp/randomcpuapp.py                -
demo-random-cpu-55954b64b4-nd89h  web-frontend  python /myapp/randomcpuapp.py                -

```
### Alternate status view
the tree flag shows the containers and pods in a tree view, with values calculated all the way up to the parent
``` shell
$ kubectl-ice status -l app=demoprobe --tree
NAMESPACE  NAME                                 READY  STARTED  RESTARTS  STATE    REASON  EXIT-CODE  SIGNAL  AGE
ice        Deployment/demo-probe                true   true     58        -        -       -          -       -
ice        └─ReplicaSet/demo-probe-765fd4d8f7   true   true     58        -        -       -          -       -
ice          └─Pod/demo-probe-765fd4d8f7-7s5d7  true   true     29        -        -       -          -       -
ice           └─Container/nginx                 true   true     0         Running  -       -          -       45h
ice           └─Container/web-frontend          true   true     29        Running  -       -          -       16m
ice          └─Pod/demo-probe-765fd4d8f7-cqr6m  true   true     29        -        -       -          -       -
ice           └─Container/nginx                 true   true     0         Running  -       -          -       45h
ice           └─Container/web-frontend          true   true     29        Running  -       -          -       16m

```
### Pick and un-mix
Using the -A flag to search all namespaces we can exclude all init containers with the --match T!=I flag. The -T flag is optional and is provided to show the init container type is not in the output
``` shell
$ kubectl-ice cpu -A -T --match T!=I
T  NAMESPACE    PODNAME                           CONTAINER                USED  REQUEST  LIMIT  %REQ      %LIMIT
C  default      web-pod                           app-watcher              5m    1m       1m     499.92    499.92
C  default      web-pod                           app-broken               0m    1m       1m     -         -
C  default      web-pod                           myapp                    5m    1m       1m     497.68    497.68
C  ice          demo-memory-7ddb58cd5b-dzp5x      web-frontend             993m  1m       1000m  99273.85  99.27
C  ice          demo-memory-7ddb58cd5b-dzp5x      nginx                    0m    1m       1000m  -         -
C  ice          demo-memory-7ddb58cd5b-hxkbt      web-frontend             997m  1m       1000m  99636.04  99.64
C  ice          demo-memory-7ddb58cd5b-hxkbt      nginx                    0m    1m       1000m  -         -
C  ice          demo-memory-7ddb58cd5b-wn7tt      web-frontend             997m  1m       1000m  99620.20  99.62
C  ice          demo-memory-7ddb58cd5b-wn7tt      nginx                    0m    1m       1000m  -         -
C  ice          demo-memory-7ddb58cd5b-xq2t4      web-frontend             994m  1m       1000m  99331.01  99.33
C  ice          demo-memory-7ddb58cd5b-xq2t4      nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-4clnc     web-frontend             129m  1m       1000m  12826.49  12.83
C  ice          demo-odd-cpu-5f947f9db4-4clnc     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-5z9w2     web-frontend             2m    1m       1000m  184.45    0.18
C  ice          demo-odd-cpu-5f947f9db4-5z9w2     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-62xjb     web-frontend             3m    1m       1000m  235.40    0.24
C  ice          demo-odd-cpu-5f947f9db4-62xjb     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-68f47     web-frontend             3m    1m       1000m  225.05    0.23
C  ice          demo-odd-cpu-5f947f9db4-68f47     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-7hlxl     web-frontend             3m    1m       1000m  223.38    0.22
C  ice          demo-odd-cpu-5f947f9db4-7hlxl     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-7r5s2     web-frontend             3m    1m       1000m  232.54    0.23
C  ice          demo-odd-cpu-5f947f9db4-7r5s2     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-8qpl5     web-frontend             3m    1m       1000m  222.98    0.22
C  ice          demo-odd-cpu-5f947f9db4-8qpl5     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-c5sv6     web-frontend             3m    1m       1000m  227.62    0.23
C  ice          demo-odd-cpu-5f947f9db4-c5sv6     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-c7scd     web-frontend             2m    1m       1000m  181.64    0.18
C  ice          demo-odd-cpu-5f947f9db4-c7scd     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-d6qz6     web-frontend             3m    1m       1000m  224.07    0.22
C  ice          demo-odd-cpu-5f947f9db4-d6qz6     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-jfwtd     web-frontend             2m    1m       1000m  196.39    0.20
C  ice          demo-odd-cpu-5f947f9db4-jfwtd     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-lkfs2     web-frontend             3m    1m       1000m  274.59    0.27
C  ice          demo-odd-cpu-5f947f9db4-lkfs2     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-q85pt     web-frontend             2m    1m       1000m  181.88    0.18
C  ice          demo-odd-cpu-5f947f9db4-q85pt     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-qdfnb     web-frontend             3m    1m       1000m  209.83    0.21
C  ice          demo-odd-cpu-5f947f9db4-qdfnb     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-tz4hj     web-frontend             3m    1m       1000m  236.00    0.24
C  ice          demo-odd-cpu-5f947f9db4-tz4hj     nginx                    0m    1m       1000m  -         -
C  ice          demo-odd-cpu-5f947f9db4-zsmwb     web-frontend             3m    1m       1000m  204.67    0.20
C  ice          demo-odd-cpu-5f947f9db4-zsmwb     nginx                    0m    1m       1000m  -         -
C  ice          demo-probe-765fd4d8f7-7s5d7       web-frontend             3m    125m     1000m  1.66      0.21
C  ice          demo-probe-765fd4d8f7-7s5d7       nginx                    0m    1m       1000m  -         -
C  ice          demo-probe-765fd4d8f7-cqr6m       web-frontend             3m    125m     1000m  2.27      0.28
C  ice          demo-probe-765fd4d8f7-cqr6m       nginx                    0m    1m       1000m  -         -
C  ice          demo-random-cpu-55954b64b4-2fnvf  web-frontend             148m  125m     1000m  118.33    14.79
C  ice          demo-random-cpu-55954b64b4-2fnvf  nginx                    0m    1m       1000m  -         -
C  ice          demo-random-cpu-55954b64b4-j5fgx  web-frontend             306m  125m     1000m  244.62    30.58
C  ice          demo-random-cpu-55954b64b4-j5fgx  nginx                    0m    1m       1000m  -         -
C  ice          demo-random-cpu-55954b64b4-kfc5q  web-frontend             260m  125m     1000m  207.35    25.92
C  ice          demo-random-cpu-55954b64b4-kfc5q  nginx                    0m    1m       1000m  -         -
C  ice          demo-random-cpu-55954b64b4-nd89h  web-frontend             521m  125m     1000m  416.55    52.07
C  ice          demo-random-cpu-55954b64b4-nd89h  nginx                    0m    1m       1000m  -         -
C  ice          web-pod                           app-watcher              5m    1m       1m     499.92    499.92
C  ice          web-pod                           app-broken               0m    1m       1m     -         -
C  ice          web-pod                           myapp                    5m    1m       1m     497.68    497.68
C  kube-system  coredns-78fcd69978-qnjtj          coredns                  1m    100m     0m     0.93      -
C  kube-system  etcd-minikube                     etcd                     9m    100m     0m     8.26      -
C  kube-system  kube-apiserver-minikube           kube-apiserver           34m   250m     0m     13.22     -
C  kube-system  kube-controller-manager-minikube  kube-controller-manager  7m    200m     0m     3.22      -
C  kube-system  kube-proxy-hdx8w                  kube-proxy               1m    0m       0m     -         -
C  kube-system  kube-scheduler-minikube           kube-scheduler           2m    100m     0m     1.37      -
C  kube-system  metrics-server-77c99ccb96-kjq9s   metrics-server           3m    100m     0m     2.40      -
C  kube-system  storage-provisioner               storage-provisioner      1m    0m       0m     -         -

```
### Filtered trees
The tree view also allows us to use the --match flag to filter based on resource type (T column) so we include deployments only providing us with a nice total of memory used for each deployment
``` shell
$ kubectl-ice mem -T --tree --match T==D
T  NAMESPACE  NAME                                      USED     REQUEST    LIMIT     %REQ     %LIMIT
D  ice        Deployment/demo-memory                    191.37Mi 11.44Mi    2334.59Mi 0.01     1.67
D  ice        Deployment/demo-odd-cpu                   119.67Mi 1556.40Mi  9338.38Mi 0.00     0.01
D  ice        Deployment/demo-probe                     6.43Mi   3.81Mi     976.56Mi  0.00     0.17
D  ice        Deployment/demo-random-cpu                42.19Mi  389.10Mi   2334.59Mi 0.00     0.01

```
### Labels and containers
you can also search specific pods and list all containers with a specific name, in this example all pods with the label app=userandomcpu are searched and only the containers that match the name web-fronteend are shown
``` shell
$ kubectl-ice cpu -l app=userandomcpu -c web-frontend
PODNAME                           CONTAINER     USED  REQUEST  LIMIT  %REQ    %LIMIT
demo-random-cpu-55954b64b4-2fnvf  web-frontend  148m  125m     1000m  118.33  14.79
demo-random-cpu-55954b64b4-j5fgx  web-frontend  306m  125m     1000m  244.62  30.58
demo-random-cpu-55954b64b4-kfc5q  web-frontend  260m  125m     1000m  207.35  25.92
demo-random-cpu-55954b64b4-nd89h  web-frontend  521m  125m     1000m  416.55  52.07

```
### Status details
using the details flag displays the timestamp and message columns
``` shell
$ kubectl-ice status -l app=myapp --details
T  PODNAME  CONTAINER    READY  STARTED  RESTARTS  STATE       REASON            EXIT-CODE  SIGNAL  TIMESTAMP            MESSAGE
I  web-pod  app-init     true   -        0         Terminated  Completed         0          0       2022-07-28 18:16:48  -
C  web-pod  app-broken   false  false    191       Waiting     CrashLoopBackOff  -          -       -                    back-off 5m0s restarting failed
C  web-pod  app-watcher  true   true     0         Running     -                 -          -       2022-07-28 18:17:07  -
C  web-pod  myapp        true   true     0         Running     -                 -          -       2022-07-28 18:17:32  -

```
### Container status
most commands work the same way including the status command which also lets you see which container(s) are causing the restarts and by using the optional --previous flag you can view the containers previous exit code
``` shell
$ kubectl-ice status -l app=myapp --previous
PODNAME  CONTAINER    STATE       REASON  EXIT-CODE  SIGNAL  TIMESTAMP            MESSAGE
web-pod  app-init     -           -       -          -       -                    -
web-pod  app-broken   Terminated  Error   1          0       2022-07-30 15:57:50  -
web-pod  app-watcher  -           -       -          -       -                    -
web-pod  myapp        -           -       -          -       -                    -

```
### Container images
need to chack on the currently configured image versions use the image command
``` shell
$ kubectl-ice image -l app=userandomcpu
PODNAME                           CONTAINER       PULL          IMAGE
demo-random-cpu-55954b64b4-2fnvf  init-myservice  IfNotPresent  busybox:1.28
demo-random-cpu-55954b64b4-2fnvf  web-frontend    Always        python:latest
demo-random-cpu-55954b64b4-2fnvf  nginx           IfNotPresent  nginx:1.7.9
demo-random-cpu-55954b64b4-j5fgx  init-myservice  IfNotPresent  busybox:1.28
demo-random-cpu-55954b64b4-j5fgx  web-frontend    Always        python:latest
demo-random-cpu-55954b64b4-j5fgx  nginx           IfNotPresent  nginx:1.7.9
demo-random-cpu-55954b64b4-kfc5q  init-myservice  IfNotPresent  busybox:1.28
demo-random-cpu-55954b64b4-kfc5q  web-frontend    Always        python:latest
demo-random-cpu-55954b64b4-kfc5q  nginx           IfNotPresent  nginx:1.7.9
demo-random-cpu-55954b64b4-nd89h  init-myservice  IfNotPresent  busybox:1.28
demo-random-cpu-55954b64b4-nd89h  web-frontend    Always        python:latest
demo-random-cpu-55954b64b4-nd89h  nginx           IfNotPresent  nginx:1.7.9

```
### Advanced labels
return cpu requests size and limits of each container where the pods have an app label that matches useoddcpu and the container name is equal to web-frontend
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-4clnc  web-frontend  129m  1m       1000m  12826.49  12.83
demo-odd-cpu-5f947f9db4-5z9w2  web-frontend  2m    1m       1000m  184.45    0.18
demo-odd-cpu-5f947f9db4-62xjb  web-frontend  3m    1m       1000m  235.40    0.24
demo-odd-cpu-5f947f9db4-68f47  web-frontend  3m    1m       1000m  225.05    0.23
demo-odd-cpu-5f947f9db4-7hlxl  web-frontend  3m    1m       1000m  223.38    0.22
demo-odd-cpu-5f947f9db4-7r5s2  web-frontend  3m    1m       1000m  232.54    0.23
demo-odd-cpu-5f947f9db4-8qpl5  web-frontend  3m    1m       1000m  222.98    0.22
demo-odd-cpu-5f947f9db4-c5sv6  web-frontend  3m    1m       1000m  227.62    0.23
demo-odd-cpu-5f947f9db4-c7scd  web-frontend  2m    1m       1000m  181.64    0.18
demo-odd-cpu-5f947f9db4-d6qz6  web-frontend  3m    1m       1000m  224.07    0.22
demo-odd-cpu-5f947f9db4-jfwtd  web-frontend  2m    1m       1000m  196.39    0.20
demo-odd-cpu-5f947f9db4-lkfs2  web-frontend  3m    1m       1000m  274.59    0.27
demo-odd-cpu-5f947f9db4-q85pt  web-frontend  2m    1m       1000m  181.88    0.18
demo-odd-cpu-5f947f9db4-qdfnb  web-frontend  3m    1m       1000m  209.83    0.21
demo-odd-cpu-5f947f9db4-tz4hj  web-frontend  3m    1m       1000m  236.00    0.24
demo-odd-cpu-5f947f9db4-zsmwb  web-frontend  3m    1m       1000m  204.67    0.20

```
### Odditites and sorting
given the listed output above the optional --oddities flag picks out the containers that have a high cpu usage when compared to the other containers listed we also sort the list in descending order by the %REQ column
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend --oddities --sort '!%REQ'
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-4clnc  web-frontend  129m  1m       1000m  12826.49  12.83

```
### Pod volumes
list all container volumes with mount points
``` shell
$ kubectl-ice volumes web-pod
CONTAINER    VOLUME                 TYPE       BACKING           SIZE  RO    MOUNT-POINT
app-init     kube-api-access-4vph2  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-watcher  app                    ConfigMap  app.py            -     false /myapp/
app-watcher  kube-api-access-4vph2  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-broken   kube-api-access-4vph2  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
myapp        app                    ConfigMap  app.py            -     false /myapp/
myapp        kube-api-access-4vph2  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount

```
### Pod exec command
retrieves the command line and any arguments specified at the container level
``` shell
$ kubectl-ice command web-pod
CONTAINER    COMMAND                   ARGUMENTS
app-init     sh -c sleep 2; exit 0     -
app-watcher  python /myapp/mainapp.py  -
app-broken   sh -c sleep 2; exit 1     -
myapp        python /myapp/mainapp.py  -

```
### Excluding rows
use the --match flag to show only the output rows where the used memory column is greater than or equal to 3MB, this has the effect of exclusing any row where the used memory column is currently under 4096kB, the value 4096 can be replaced with any whole number in kilobytes
``` shell
$ kubectl-ice mem -l app=userandomcpu --match 'used>=4096'
PODNAME                           CONTAINER       USED    REQUEST  LIMIT  %REQ    %LIMIT
demo-random-cpu-55954b64b4-2fnvf  web-frontend    8.29Mi  1M       256M   868.76  3.39
demo-random-cpu-55954b64b4-j5fgx  web-frontend    8.40Mi  1M       256M   881.05  3.44
demo-random-cpu-55954b64b4-kfc5q  web-frontend    8.36Mi  1M       256M   876.13  3.42
demo-random-cpu-55954b64b4-nd89h  web-frontend    8.36Mi  1M       256M   876.54  3.42

```
### Extra selections
using the --select flag allows you to filter the pod selection to only pods that have a priorityClassName thats equal to system-cluster-critical, you can also match against priority
``` shell
$ kubectl-ice status --select 'priorityClassName=system-cluster-critical' -A
NAMESPACE    PODNAME                          CONTAINER       READY  STARTED  RESTARTS  STATE    REASON  EXIT-CODE  SIGNAL  AGE
kube-system  coredns-78fcd69978-qnjtj         coredns         true   true     0         Running  -       -          -       17d
kube-system  metrics-server-77c99ccb96-kjq9s  metrics-server  true   true     0         Running  -       -          -       17d

```
### Security information
listing runAsUser and runAsGroup settings along with other related container security information
``` shell
$ kubectl-ice security -n kube-system
PODNAME                           CONTAINER                ALLOW_PRIVILEGE_ESCALATION  PRIVILEGED  RO_ROOT_FS  RUN_AS_NON_ROOT  RUN_AS_USER  RUN_AS_GROUP
coredns-78fcd69978-qnjtj          coredns                  false                       -           true        -                -            -
etcd-minikube                     etcd                     -                           -           -           -                -            -
kube-apiserver-minikube           kube-apiserver           -                           -           -           -                -            -
kube-controller-manager-minikube  kube-controller-manager  -                           -           -           -                -            -
kube-proxy-hdx8w                  kube-proxy               -                           true        -           -                -            -
kube-scheduler-minikube           kube-scheduler           -                           -           -           -                -            -
metrics-server-77c99ccb96-kjq9s   metrics-server           -                           -           true        true             1000         -
storage-provisioner               storage-provisioner      -                           -           -           -                -            -

```
### POSIX capabilities
display configured capabilities related to each container
``` shell
$ kubectl-ice capabilities -n kube-system
PODNAME                           CONTAINER                ADD               DROP
coredns-78fcd69978-qnjtj          coredns                  NET_BIND_SERVICE  all
etcd-minikube                     etcd                     -                 -
kube-apiserver-minikube           kube-apiserver           -                 -
kube-controller-manager-minikube  kube-controller-manager  -                 -
kube-proxy-hdx8w                  kube-proxy               -                 -
kube-scheduler-minikube           kube-scheduler           -                 -
metrics-server-77c99ccb96-kjq9s   metrics-server           -                 -
storage-provisioner               storage-provisioner      -                 -

```
### Column labels
with the --node-label and --pod-label flags its possible to show the values of the labels as columns in the output table
``` shell
$ kubectl-ice status --node-label "beta.kubernetes.io/os" --pod-label "component" -n kube-system
PODNAME                           CONTAINER                beta.kubernetes.io/os  component                READY  STARTED  RESTARTS  STATE    REASON  EXIT-CODE  SIGNAL  AGE
coredns-78fcd69978-qnjtj          coredns                  linux                  -                        true   true     0         Running  -       -          -       17d
etcd-minikube                     etcd                     linux                  etcd                     true   true     0         Running  -       -          -       17d
kube-apiserver-minikube           kube-apiserver           linux                  kube-apiserver           true   true     0         Running  -       -          -       17d
kube-controller-manager-minikube  kube-controller-manager  linux                  kube-controller-manager  true   true     0         Running  -       -          -       17d
kube-proxy-hdx8w                  kube-proxy               linux                  -                        true   true     0         Running  -       -          -       17d
kube-scheduler-minikube           kube-scheduler           linux                  kube-scheduler           true   true     0         Running  -       -          -       17d
metrics-server-77c99ccb96-kjq9s   metrics-server           linux                  -                        true   true     0         Running  -       -          -       17d
storage-provisioner               storage-provisioner      linux                  -                        true   true     1         Running  -       -          -       17d

```
