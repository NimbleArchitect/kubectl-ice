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
### Named containers
the optional container flag (-c) searchs all selected pods and lists only containers that match the name web-frontend
``` shell
$ kubectl-ice command -c web-frontend
PODNAME                           CONTAINER     COMMAND                                      ARGUMENTS
demo-memory-7ddb58cd5b-5kf9g      web-frontend  python /myapp/halfmemapp.py                  -
demo-memory-7ddb58cd5b-csbds      web-frontend  python /myapp/halfmemapp.py                  -
demo-memory-7ddb58cd5b-d4zwp      web-frontend  python /myapp/halfmemapp.py                  -
demo-memory-7ddb58cd5b-pdm9c      web-frontend  python /myapp/halfmemapp.py                  -
demo-odd-cpu-5f947f9db4-2g7p2     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-59jm6     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-6gzw7     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-6s97l     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-86mb8     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-cwvdq     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-dcg8p     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-fhs8q     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-gzcrm     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-hf872     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-hft68     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-jp8fw     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-k2gtp     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-kj8s7     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-qtxp2     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-vg2d5     web-frontend  python /myapp/oddcpuapp.py                   -
demo-probe-765fd4d8f7-n6kc7       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
demo-probe-765fd4d8f7-x2zr6       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
demo-random-cpu-55954b64b4-9t7m2  web-frontend  python /myapp/randomcpuapp.py                -
demo-random-cpu-55954b64b4-km6bg  web-frontend  python /myapp/randomcpuapp.py                -
demo-random-cpu-55954b64b4-knc6n  web-frontend  python /myapp/randomcpuapp.py                -
demo-random-cpu-55954b64b4-vr4hg  web-frontend  python /myapp/randomcpuapp.py                -

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
### Labels and containers
you can also search specific pods and list all containers with a specific name, in this example all pods with the label app=userandomcpu are searched and only the containers that match the name web-fronteend are shown
``` shell
$ kubectl-ice cpu -l app=userandomcpu -c web-frontend
PODNAME                           CONTAINER     USED  REQUEST  LIMIT  %REQ    %LIMIT
demo-random-cpu-55954b64b4-9t7m2  web-frontend  569m  125m     1000m  454.43  56.80
demo-random-cpu-55954b64b4-km6bg  web-frontend  449m  125m     1000m  358.72  44.84
demo-random-cpu-55954b64b4-knc6n  web-frontend  206m  125m     1000m  164.19  20.52
demo-random-cpu-55954b64b4-vr4hg  web-frontend  456m  125m     1000m  364.08  45.51

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
### Container images
need to chack on the currently configured image versions use the image command
``` shell
$ kubectl-ice image -l app=userandomcpu
PODNAME                           CONTAINER       PULL          IMAGE
demo-random-cpu-55954b64b4-9t7m2  init-myservice  IfNotPresent  busybox:1.28
demo-random-cpu-55954b64b4-9t7m2  web-frontend    Always        python:latest
demo-random-cpu-55954b64b4-9t7m2  nginx           IfNotPresent  nginx:1.7.9
demo-random-cpu-55954b64b4-km6bg  init-myservice  IfNotPresent  busybox:1.28
demo-random-cpu-55954b64b4-km6bg  web-frontend    Always        python:latest
demo-random-cpu-55954b64b4-km6bg  nginx           IfNotPresent  nginx:1.7.9
demo-random-cpu-55954b64b4-knc6n  init-myservice  IfNotPresent  busybox:1.28
demo-random-cpu-55954b64b4-knc6n  web-frontend    Always        python:latest
demo-random-cpu-55954b64b4-knc6n  nginx           IfNotPresent  nginx:1.7.9
demo-random-cpu-55954b64b4-vr4hg  init-myservice  IfNotPresent  busybox:1.28
demo-random-cpu-55954b64b4-vr4hg  web-frontend    Always        python:latest
demo-random-cpu-55954b64b4-vr4hg  nginx           IfNotPresent  nginx:1.7.9

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
### Pod volumes
list all container volumes with mount points
``` shell
$ kubectl-ice volumes web-pod
CONTAINER    VOLUME                 TYPE       BACKING           SIZE  RO    MOUNT-POINT
app-init     kube-api-access-c47d7  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-watcher  app                    ConfigMap  app.py            -     false /myapp/
app-watcher  kube-api-access-c47d7  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-broken   kube-api-access-c47d7  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
myapp        app                    ConfigMap  app.py            -     false /myapp/
myapp        kube-api-access-c47d7  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount

```
### Pod exec command
retrieves the command line and any arguments specified at the container level
``` shell
$ kubectl-ice command web-pod
CONTAINER       COMMAND                   ARGUMENTS
app-init        sh -c sleep 2; exit 0     -
app-watcher     python /myapp/mainapp.py  -
app-broken      sh -c sleep 2; exit 1     -
myapp           python /myapp/mainapp.py  -
debugger-k5znj  -                         -

```
### Excluding rows
use the --match flag to show only the output rows where the used memory column is greater than or equal to 3MB, this has the effect of exclusing any row where the used memory column is currently under 4096kB, the value 4096 can be replaced with any whole number in kilobytes
``` shell
$ kubectl-ice mem -l app=userandomcpu --match 'used>=4096'
PODNAME                           CONTAINER     USED    REQUEST  LIMIT  %REQ    %LIMIT
demo-random-cpu-55954b64b4-9t7m2  web-frontend  6.71Mi  1M       256M   704.10  2.75
demo-random-cpu-55954b64b4-km6bg  web-frontend  8.29Mi  1M       256M   869.17  3.40
demo-random-cpu-55954b64b4-knc6n  web-frontend  5.41Mi  1M       256M   566.89  2.21
demo-random-cpu-55954b64b4-vr4hg  web-frontend  8.63Mi  1M       256M   904.81  3.53

```
### Extra selections
using the --select flag allows you to filter the pod selection to only pods that have a priorityClassName thats equal to system-cluster-critical, you can also match against priority
``` shell
$ kubectl-ice status --select 'priorityClassName=system-cluster-critical' -A
NAMESPACE    PODNAME                          CONTAINER       READY  STARTED  RESTARTS  STATE    REASON  EXIT-CODE  SIGNAL  AGE
kube-system  coredns-78fcd69978-qnjtj         coredns         true   true     1         Running  -       -          -       46h
kube-system  metrics-server-77c99ccb96-kjq9s  metrics-server  true   true     1         Running  -       -          -       46h

```
### Security information
listing runAsUser and runAsGroup settings along with other related container security information
``` shell
$ kubectl-ice security -n kube-system
PODNAME                           CONTAINER                ALLOW_PRIVILEGE_ESCALATION  PRIVILEGED  RO_ROOT_FS  RUN_AS_NON_ROOT  RUN_AS_USER  RUN_AS_GROUP
coredns-78fcd69978-qnjtj          coredns                  false                       -           true        -                -            -
etcd-minikube                     etcd                     -                           -           -           -                -            -
kindnet-2mhhp                     kindnet-cni              -                           false       -           -                -            -
kindnet-qpdhv                     kindnet-cni              -                           false       -           -                -            -
kube-apiserver-minikube           kube-apiserver           -                           -           -           -                -            -
kube-controller-manager-minikube  kube-controller-manager  -                           -           -           -                -            -
kube-proxy-f98q8                  kube-proxy               -                           true        -           -                -            -
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
kindnet-2mhhp                     kindnet-cni              NET_RAW,NET_ADMIN -
kindnet-qpdhv                     kindnet-cni              NET_RAW,NET_ADMIN -
kube-apiserver-minikube           kube-apiserver           -                 -
kube-controller-manager-minikube  kube-controller-manager  -                 -
kube-proxy-f98q8                  kube-proxy               -                 -
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
coredns-78fcd69978-qnjtj          coredns                  linux                  -                        true   true     1         Running  -       -          -       46h
etcd-minikube                     etcd                     linux                  etcd                     true   true     1         Running  -       -          -       46h
kindnet-2mhhp                     kindnet-cni              linux                  -                        true   true     1         Running  -       -          -       46h
kindnet-qpdhv                     kindnet-cni              linux                  -                        true   true     1         Running  -       -          -       46h
kube-apiserver-minikube           kube-apiserver           linux                  kube-apiserver           true   true     1         Running  -       -          -       46h
kube-controller-manager-minikube  kube-controller-manager  linux                  kube-controller-manager  true   true     1         Running  -       -          -       46h
kube-proxy-f98q8                  kube-proxy               linux                  -                        true   true     1         Running  -       -          -       46h
kube-proxy-hdx8w                  kube-proxy               linux                  -                        true   true     1         Running  -       -          -       46h
kube-scheduler-minikube           kube-scheduler           linux                  kube-scheduler           true   true     1         Running  -       -          -       46h
metrics-server-77c99ccb96-kjq9s   metrics-server           linux                  -                        true   true     1         Running  -       -          -       46h
storage-provisioner               storage-provisioner      linux                  -                        true   true     2         Running  -       -          -       46h

```
