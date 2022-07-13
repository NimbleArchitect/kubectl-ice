### Single pod info
Shows the currently used memory along with the configured memory requests and limits of all containers (side cars) in the pod named web-pod
``` shell
$ kubectl-ice memory web-pod
CONTAINER    USED    REQUEST  LIMIT  %REQ    %LIMIT
app-init     0       0        0      -       -
app-watcher  6.52Mi  1M       512M   683.21  1.33
app-broken   0       1M       512M   -       -
myapp        6.79Mi  1M       256M   711.48  2.78

```
### Using labels
using labels you can search all pods that are part of a deployment where the label app matches demoprobe and list selected information about the containers in each pod, this example shows the currently configured probe information and gives details of configured startup, readiness and liveness probes of each container
``` shell
$ kubectl-ice probes -l app=demoprobe
PODNAME                      CONTAINER     PROBE     DELAY  PERIOD  TIMEOUT  SUCCESS  FAILURE  CHECK    ACTION
demo-probe-765fd4d8f7-r5rq7  web-frontend  liveness  10     5       1        1        3        Exec     /bin/true
demo-probe-765fd4d8f7-r5rq7  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-765fd4d8f7-r5rq7  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/
demo-probe-765fd4d8f7-wxt7f  web-frontend  liveness  10     5       1        1        3        Exec     /bin/true
demo-probe-765fd4d8f7-wxt7f  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-765fd4d8f7-wxt7f  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/

```
### Named containers
the optional container flag (-c) searchs all selected pods and lists only containers that match the name web-frontend
``` shell
$ kubectl-ice command -c web-frontend
PODNAME                           CONTAINER     COMMAND                                      ARGUMENTS
demo-memory-7ddb58cd5b-g76gd      web-frontend  python /myapp/halfmemapp.py                  -
demo-memory-7ddb58cd5b-ngbng      web-frontend  python /myapp/halfmemapp.py                  -
demo-memory-7ddb58cd5b-r22kp      web-frontend  python /myapp/halfmemapp.py                  -
demo-memory-7ddb58cd5b-r6bgn      web-frontend  python /myapp/halfmemapp.py                  -
demo-odd-cpu-5f947f9db4-2wgdp     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-56hnq     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-77q9n     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-9kwbm     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-9t7vg     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-g2258     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-hd6ld     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-mnznh     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-n9zvc     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-qtxmf     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-ssv2x     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-v24n4     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-wfg78     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-x5k44     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-xbrhk     web-frontend  python /myapp/oddcpuapp.py                   -
demo-odd-cpu-5f947f9db4-zt5vw     web-frontend  python /myapp/oddcpuapp.py                   -
demo-probe-765fd4d8f7-r5rq7       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
demo-probe-765fd4d8f7-wxt7f       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
demo-random-cpu-55954b64b4-l4b44  web-frontend  python /myapp/randomcpuapp.py                -
demo-random-cpu-55954b64b4-nwm9c  web-frontend  python /myapp/randomcpuapp.py                -
demo-random-cpu-55954b64b4-vnrk5  web-frontend  python /myapp/randomcpuapp.py                -
demo-random-cpu-55954b64b4-xnxhd  web-frontend  python /myapp/randomcpuapp.py                -

```
### Alternate status view
the tree flag shows the containers and pods in a tree like view
``` shell
$ kubectl-ice status -l app=myapp --tree
NAMESPACE  NAME                      READY  STARTED  RESTARTS  STATE       REASON            EXIT-CODE  SIGNAL  AGE
ice        Pod/web-pod               -      -        0         Running     -                 -          -       104m
ice        └─InitContainer/app-init  true   -        0         Terminated  Completed         0          0       104m
ice        └─Container/app-broken    false  false    24        Waiting     CrashLoopBackOff  -          -       -
ice        └─Container/app-watcher   true   true     0         Running     -                 -          -       103m
ice        └─Container/myapp         true   true     0         Running     -                 -          -       103m

```
### Pick and un-mix
Using the -A flag to search all namespaces we can exclude all init containers with the --match T!=I flag. The -T flag is optional and is provided to show the init container type is not in the output
``` shell
$ kubectl-ice cpu -A -T --match T!=I
T  NAMESPACE    PODNAME                           CONTAINER                USED  REQUEST  LIMIT  %REQ      %LIMIT
S  default      web-pod                           app-watcher              5m    1m       1m     495.46    495.46
S  default      web-pod                           app-broken               0m    1m       1m     -         -
S  default      web-pod                           myapp                    6m    1m       1m     516.79    516.79
S  ice          demo-memory-7ddb58cd5b-g76gd      web-frontend             997m  1m       1000m  99633.77  99.63
S  ice          demo-memory-7ddb58cd5b-g76gd      nginx                    0m    1m       1000m  -         -
S  ice          demo-memory-7ddb58cd5b-ngbng      web-frontend             996m  1m       1000m  99528.84  99.53
S  ice          demo-memory-7ddb58cd5b-ngbng      nginx                    0m    1m       1000m  -         -
S  ice          demo-memory-7ddb58cd5b-r22kp      web-frontend             997m  1m       1000m  99671.11  99.67
S  ice          demo-memory-7ddb58cd5b-r22kp      nginx                    0m    1m       1000m  -         -
S  ice          demo-memory-7ddb58cd5b-r6bgn      web-frontend             997m  1m       1000m  99683.65  99.68
S  ice          demo-memory-7ddb58cd5b-r6bgn      nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-2wgdp     web-frontend             3m    1m       1000m  276.45    0.28
S  ice          demo-odd-cpu-5f947f9db4-2wgdp     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-56hnq     web-frontend             2m    1m       1000m  199.89    0.20
S  ice          demo-odd-cpu-5f947f9db4-56hnq     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-77q9n     web-frontend             123m  1m       1000m  12287.04  12.29
S  ice          demo-odd-cpu-5f947f9db4-77q9n     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-9kwbm     web-frontend             3m    1m       1000m  234.97    0.23
S  ice          demo-odd-cpu-5f947f9db4-9kwbm     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-9t7vg     web-frontend             3m    1m       1000m  209.26    0.21
S  ice          demo-odd-cpu-5f947f9db4-9t7vg     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-g2258     web-frontend             3m    1m       1000m  271.59    0.27
S  ice          demo-odd-cpu-5f947f9db4-g2258     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-hd6ld     web-frontend             3m    1m       1000m  204.26    0.20
S  ice          demo-odd-cpu-5f947f9db4-hd6ld     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-mnznh     web-frontend             113m  1m       1000m  11287.89  11.29
S  ice          demo-odd-cpu-5f947f9db4-mnznh     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-n9zvc     web-frontend             3m    1m       1000m  215.90    0.22
S  ice          demo-odd-cpu-5f947f9db4-n9zvc     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-qtxmf     web-frontend             3m    1m       1000m  238.15    0.24
S  ice          demo-odd-cpu-5f947f9db4-qtxmf     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-ssv2x     web-frontend             3m    1m       1000m  253.31    0.25
S  ice          demo-odd-cpu-5f947f9db4-ssv2x     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-v24n4     web-frontend             3m    1m       1000m  261.50    0.26
S  ice          demo-odd-cpu-5f947f9db4-v24n4     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-wfg78     web-frontend             3m    1m       1000m  240.82    0.24
S  ice          demo-odd-cpu-5f947f9db4-wfg78     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-x5k44     web-frontend             2m    1m       1000m  198.79    0.20
S  ice          demo-odd-cpu-5f947f9db4-x5k44     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-xbrhk     web-frontend             3m    1m       1000m  254.32    0.25
S  ice          demo-odd-cpu-5f947f9db4-xbrhk     nginx                    0m    1m       1000m  -         -
S  ice          demo-odd-cpu-5f947f9db4-zt5vw     web-frontend             3m    1m       1000m  229.52    0.23
S  ice          demo-odd-cpu-5f947f9db4-zt5vw     nginx                    0m    1m       1000m  -         -
S  ice          demo-probe-765fd4d8f7-r5rq7       web-frontend             3m    125m     1000m  2.02      0.25
S  ice          demo-probe-765fd4d8f7-r5rq7       nginx                    1m    1m       1000m  1.24      0.00
S  ice          demo-probe-765fd4d8f7-wxt7f       web-frontend             3m    125m     1000m  2.05      0.26
S  ice          demo-probe-765fd4d8f7-wxt7f       nginx                    0m    1m       1000m  -         -
S  ice          demo-random-cpu-55954b64b4-l4b44  web-frontend             145m  125m     1000m  115.28    14.41
S  ice          demo-random-cpu-55954b64b4-l4b44  nginx                    0m    1m       1000m  -         -
S  ice          demo-random-cpu-55954b64b4-nwm9c  web-frontend             249m  125m     1000m  199.01    24.88
S  ice          demo-random-cpu-55954b64b4-nwm9c  nginx                    1m    1m       1000m  1.73      0.00
S  ice          demo-random-cpu-55954b64b4-vnrk5  web-frontend             279m  125m     1000m  222.59    27.82
S  ice          demo-random-cpu-55954b64b4-vnrk5  nginx                    1m    1m       1000m  0.77      0.00
S  ice          demo-random-cpu-55954b64b4-xnxhd  web-frontend             254m  125m     1000m  202.81    25.35
S  ice          demo-random-cpu-55954b64b4-xnxhd  nginx                    1m    1m       1000m  1.10      0.00
S  ice          web-pod                           app-watcher              5m    1m       1m     495.46    495.46
S  ice          web-pod                           app-broken               0m    1m       1m     -         -
S  ice          web-pod                           myapp                    6m    1m       1m     516.79    516.79
S  kube-system  coredns-78fcd69978-qnjtj          coredns                  1m    100m     0m     0.88      -
S  kube-system  etcd-minikube                     etcd                     8m    100m     0m     7.16      -
S  kube-system  kube-apiserver-minikube           kube-apiserver           29m   250m     0m     11.42     -
S  kube-system  kube-controller-manager-minikube  kube-controller-manager  11m   200m     0m     5.29      -
S  kube-system  kube-proxy-hdx8w                  kube-proxy               1m    0m       0m     -         -
S  kube-system  kube-scheduler-minikube           kube-scheduler           2m    100m     0m     1.06      -
S  kube-system  metrics-server-77c99ccb96-kjq9s   metrics-server           2m    100m     0m     1.78      -
S  kube-system  storage-provisioner               storage-provisioner      1m    0m       0m     -         -

```
### Labels and containers
you can also search specific pods and list all containers with a specific name, in this example all pods with the label app=userandomcpu are searched and only the containers that match the name web-fronteend are shown
``` shell
$ kubectl-ice cpu -l app=userandomcpu -c web-frontend
PODNAME                           CONTAINER     USED  REQUEST  LIMIT  %REQ    %LIMIT
demo-random-cpu-55954b64b4-l4b44  web-frontend  145m  125m     1000m  115.28  14.41
demo-random-cpu-55954b64b4-nwm9c  web-frontend  249m  125m     1000m  199.01  24.88
demo-random-cpu-55954b64b4-vnrk5  web-frontend  279m  125m     1000m  222.59  27.82
demo-random-cpu-55954b64b4-xnxhd  web-frontend  254m  125m     1000m  202.81  25.35

```
### Status details
using the details flag displays the timestamp and message columns
``` shell
$ kubectl-ice status -l app=myapp --details
T  PODNAME  CONTAINER    READY  STARTED  RESTARTS  STATE       REASON            EXIT-CODE  SIGNAL  TIMESTAMP            MESSAGE
I  web-pod  app-init     true   -        0         Terminated  Completed         0          0       2022-07-13 16:24:31  -
S  web-pod  app-broken   false  false    24        Waiting     CrashLoopBackOff  -          -       -                    back-off 5m0s restarting failed
S  web-pod  app-watcher  true   true     0         Running     -                 -          -       2022-07-13 16:25:04  -
S  web-pod  myapp        true   true     0         Running     -                 -          -       2022-07-13 16:25:17  -

```
### Container status
most commands work the same way including the status command which also lets you see which container(s) are causing the restarts and by using the optional --previous flag you can view the containers previous exit code
``` shell
$ kubectl-ice status -l app=myapp --previous
PODNAME  CONTAINER    STATE       REASON  EXIT-CODE  SIGNAL  TIMESTAMP            MESSAGE
web-pod  app-init     -           -       -          -       -                    -
web-pod  app-broken   Terminated  Error   1          0       2022-07-13 18:07:49  -
web-pod  app-watcher  -           -       -          -       -                    -
web-pod  myapp        -           -       -          -       -                    -

```
### Container images
need to chack on the currently configured image versions use the image command
``` shell
$ kubectl-ice image -l app=userandomcpu
PODNAME                           CONTAINER       PULL          IMAGE
demo-random-cpu-55954b64b4-l4b44  init-myservice  IfNotPresent  busybox:1.28
demo-random-cpu-55954b64b4-l4b44  web-frontend    Always        python:latest
demo-random-cpu-55954b64b4-l4b44  nginx           IfNotPresent  nginx:1.7.9
demo-random-cpu-55954b64b4-nwm9c  init-myservice  IfNotPresent  busybox:1.28
demo-random-cpu-55954b64b4-nwm9c  web-frontend    Always        python:latest
demo-random-cpu-55954b64b4-nwm9c  nginx           IfNotPresent  nginx:1.7.9
demo-random-cpu-55954b64b4-vnrk5  init-myservice  IfNotPresent  busybox:1.28
demo-random-cpu-55954b64b4-vnrk5  web-frontend    Always        python:latest
demo-random-cpu-55954b64b4-vnrk5  nginx           IfNotPresent  nginx:1.7.9
demo-random-cpu-55954b64b4-xnxhd  init-myservice  IfNotPresent  busybox:1.28
demo-random-cpu-55954b64b4-xnxhd  web-frontend    Always        python:latest
demo-random-cpu-55954b64b4-xnxhd  nginx           IfNotPresent  nginx:1.7.9

```
### Advanced labels
return memory requests size and limits of each container where the pods have an app label that matches useoddcpu and the container name is equal to web-frontend
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-2wgdp  web-frontend  3m    1m       1000m  276.45    0.28
demo-odd-cpu-5f947f9db4-56hnq  web-frontend  2m    1m       1000m  199.89    0.20
demo-odd-cpu-5f947f9db4-77q9n  web-frontend  123m  1m       1000m  12287.04  12.29
demo-odd-cpu-5f947f9db4-9kwbm  web-frontend  3m    1m       1000m  234.97    0.23
demo-odd-cpu-5f947f9db4-9t7vg  web-frontend  3m    1m       1000m  209.26    0.21
demo-odd-cpu-5f947f9db4-g2258  web-frontend  3m    1m       1000m  271.59    0.27
demo-odd-cpu-5f947f9db4-hd6ld  web-frontend  3m    1m       1000m  204.26    0.20
demo-odd-cpu-5f947f9db4-mnznh  web-frontend  113m  1m       1000m  11287.89  11.29
demo-odd-cpu-5f947f9db4-n9zvc  web-frontend  3m    1m       1000m  215.90    0.22
demo-odd-cpu-5f947f9db4-qtxmf  web-frontend  3m    1m       1000m  238.15    0.24
demo-odd-cpu-5f947f9db4-ssv2x  web-frontend  3m    1m       1000m  253.31    0.25
demo-odd-cpu-5f947f9db4-v24n4  web-frontend  3m    1m       1000m  261.50    0.26
demo-odd-cpu-5f947f9db4-wfg78  web-frontend  3m    1m       1000m  240.82    0.24
demo-odd-cpu-5f947f9db4-x5k44  web-frontend  2m    1m       1000m  198.79    0.20
demo-odd-cpu-5f947f9db4-xbrhk  web-frontend  3m    1m       1000m  254.32    0.25
demo-odd-cpu-5f947f9db4-zt5vw  web-frontend  3m    1m       1000m  229.52    0.23

```
### Odditites and sorting
given the listed output above the optional --oddities flag picks out the containers that have a high cpu usage when compared to the other containers listed we also sort the list in descending order by the %REQ column
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend --oddities --sort '!%REQ'
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-77q9n  web-frontend  123m  1m       1000m  12287.04  12.29
demo-odd-cpu-5f947f9db4-mnznh  web-frontend  113m  1m       1000m  11287.89  11.29

```
### Pod volumes
list all container volumes with mount points
``` shell
$ kubectl-ice volumes web-pod
CONTAINER    VOLUME                 TYPE       BACKING           SIZE  RO    MOUNT-POINT
app-init     kube-api-access-7cx8b  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-watcher  app                    ConfigMap  app.py            -     false /myapp/
app-watcher  kube-api-access-7cx8b  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-broken   kube-api-access-7cx8b  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
myapp        app                    ConfigMap  app.py            -     false /myapp/
myapp        kube-api-access-7cx8b  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount

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
demo-random-cpu-55954b64b4-l4b44  web-frontend    9.25Mi  1M       256M   969.52  3.79
demo-random-cpu-55954b64b4-l4b44  nginx           4.19Mi  1M       256M   439.50  1.72
demo-random-cpu-55954b64b4-nwm9c  web-frontend    9.05Mi  1M       256M   949.45  3.71
demo-random-cpu-55954b64b4-vnrk5  web-frontend    9.08Mi  1M       256M   951.91  3.72
demo-random-cpu-55954b64b4-xnxhd  web-frontend    8.82Mi  1M       256M   924.88  3.61

```
### Extra selections
using the --select flag allows you to filter the pod selection to only pods that have a priorityClassName thats equal to system-cluster-critical, you can also match against priority
``` shell
$ kubectl-ice status --select 'priorityClassName=system-cluster-critical' -A
NAMESPACE    PODNAME                          CONTAINER       READY  STARTED  RESTARTS  STATE    REASON  EXIT-CODE  SIGNAL  AGE
kube-system  coredns-78fcd69978-qnjtj         coredns         true   true     0         Running  -       -          -       9h
kube-system  metrics-server-77c99ccb96-kjq9s  metrics-server  true   true     0         Running  -       -          -       8h

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
coredns-78fcd69978-qnjtj          coredns                  linux                  -                        true   true     0         Running  -       -          -       9h
etcd-minikube                     etcd                     linux                  etcd                     true   true     0         Running  -       -          -       9h
kube-apiserver-minikube           kube-apiserver           linux                  kube-apiserver           true   true     0         Running  -       -          -       9h
kube-controller-manager-minikube  kube-controller-manager  linux                  kube-controller-manager  true   true     0         Running  -       -          -       9h
kube-proxy-hdx8w                  kube-proxy               linux                  -                        true   true     0         Running  -       -          -       9h
kube-scheduler-minikube           kube-scheduler           linux                  kube-scheduler           true   true     0         Running  -       -          -       9h
metrics-server-77c99ccb96-kjq9s   metrics-server           linux                  -                        true   true     0         Running  -       -          -       8h
storage-provisioner               storage-provisioner      linux                  -                        true   true     1         Running  -       -          -       9h

```
