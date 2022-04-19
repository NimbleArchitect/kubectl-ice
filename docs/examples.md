### Single pod info
Shows the currently used memory along with the configured memory requests and limits of all containers (side cars) in the pod named web-pod
``` shell
$ kubectl-ice memory web-pod
T  CONTAINER    USED    REQUEST  LIMIT  %REQ    %LIMIT
I  app-init     0       0        0      -       -
S  app-watcher  0       1M       512M   -       -
S  app-broken   0       1M       512M   -       -
S  myapp        0.01Gi  1M       256M   550.09  2.15

```
### Using labels
using labels you can search all pods that are part of a deployment where the label app matches demoprobe and list selected information about the containers in each pod, this example shows the currently configured probe information and gives details of configured startup, readiness and liveness probes of each container
``` shell
$ kubectl-ice probes -l app=demoprobe
PODNAME                      CONTAINER     PROBE     DELAY  PERIOD  TIMEOUT  SUCCESS  FAILURE  CHECK    ACTION
demo-probe-76b66d5766-ckhmt  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-ckhmt  web-frontend  liveness  5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-ckhmt  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/
demo-probe-76b66d5766-zk8pp  web-frontend  liveness  5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-zk8pp  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-zk8pp  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/

```
### Named containers
the optional container flag (-c) searchs all selected pods and lists only containers that match the name web-frontend
``` shell
$ kubectl-ice command -c web-frontend
T  PODNAME                           CONTAINER     COMMAND                                      ARGUMENTS
S  demo-memory-7ddb58cd5b-lj7xb      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-nhd8c      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-pnk9r      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-pqxsg      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-odd-cpu-5f947f9db4-4g2gk     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-8pkpq     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-9rdh4     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-b227c     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-dldv5     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-fg9wj     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-jnfck     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-kktmw     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-mgbpv     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-p9mxv     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-r5rk4     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-rhktb     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-tthwv     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-vqnm9     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-w8w2z     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-x65bq     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-probe-76b66d5766-ckhmt       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
S  demo-probe-76b66d5766-zk8pp       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
S  demo-random-cpu-669b7888b9-46tzm  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-56hbd  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-jl4vg  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-rbs6d  web-frontend  python /myapp/randomcpuapp.py                -

```
### Labels and containers
you can also search specific pods and list all containers with a specific name, in this example all pods with the label app=userandomcpu are searched and only the containers that match the name web-fronteend are shown
``` shell
$ kubectl-ice cpu -l app=userandomcpu -c web-frontend
T  PODNAME                           CONTAINER     USED  REQUEST  LIMIT  %REQ  %LIMIT
S  demo-random-cpu-669b7888b9-46tzm  web-frontend  0     125      1000   -     -
S  demo-random-cpu-669b7888b9-56hbd  web-frontend  0     125      1000   -     -
S  demo-random-cpu-669b7888b9-jl4vg  web-frontend  0     125      1000   -     -
S  demo-random-cpu-669b7888b9-rbs6d  web-frontend  0     125      1000   -     -

```
### Container status
most commands work the same way including the status command which also lets you see which container(s) are causing the restarts and by using the optional --previous flag you can view the containers previous exit code
``` shell
$ kubectl-ice status -l app=myapp --previous
T  PODNAME  CONTAINER    STATE       REASON  EXIT-CODE  SIGNAL  TIMESTAMP                      MESSAGE
S  web-pod  app-broken   Terminated  Error   1          0       2022-04-19 16:03:03 +0100 BST  -
S  web-pod  app-watcher  Terminated  Error   2          0       2022-04-19 15:59:18 +0100 BST  -
S  web-pod  myapp        -           -       -          -       -                              -
I  web-pod  app-init     -           -       -          -       -                              -

```
### Container images
need to chack on the currently configured image versions use the image command
``` shell
$ kubectl-ice image -l app=userandomcpu
T  PODNAME                           CONTAINER       PULL          IMAGE
S  demo-random-cpu-669b7888b9-46tzm  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-46tzm  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-46tzm  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-56hbd  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-56hbd  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-56hbd  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-jl4vg  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-jl4vg  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-jl4vg  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-rbs6d  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-rbs6d  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-rbs6d  init-myservice  IfNotPresent  busybox:1.28

```
### Advanced labels
return memory requests size and limits of each container where the pods have an app label that matches useoddcpu and the container name is equal to nginx
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend
T  PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
S  demo-odd-cpu-5f947f9db4-4g2gk  web-frontend  2     1        1000   141.36    0.14
S  demo-odd-cpu-5f947f9db4-8pkpq  web-frontend  2     1        1000   132.28    0.13
S  demo-odd-cpu-5f947f9db4-9rdh4  web-frontend  2     1        1000   147.10    0.15
S  demo-odd-cpu-5f947f9db4-b227c  web-frontend  2     1        1000   157.30    0.16
S  demo-odd-cpu-5f947f9db4-dldv5  web-frontend  90    1        1000   8947.00   8.95
S  demo-odd-cpu-5f947f9db4-fg9wj  web-frontend  107   1        1000   10640.05  10.64
S  demo-odd-cpu-5f947f9db4-jnfck  web-frontend  2     1        1000   147.75    0.15
S  demo-odd-cpu-5f947f9db4-kktmw  web-frontend  2     1        1000   149.83    0.15
S  demo-odd-cpu-5f947f9db4-mgbpv  web-frontend  2     1        1000   138.21    0.14
S  demo-odd-cpu-5f947f9db4-p9mxv  web-frontend  2     1        1000   147.88    0.15
S  demo-odd-cpu-5f947f9db4-r5rk4  web-frontend  2     1        1000   140.73    0.14
S  demo-odd-cpu-5f947f9db4-rhktb  web-frontend  2     1        1000   145.02    0.15
S  demo-odd-cpu-5f947f9db4-tthwv  web-frontend  105   1        1000   10426.41  10.43
S  demo-odd-cpu-5f947f9db4-vqnm9  web-frontend  2     1        1000   131.28    0.13
S  demo-odd-cpu-5f947f9db4-w8w2z  web-frontend  2     1        1000   134.02    0.13
S  demo-odd-cpu-5f947f9db4-x65bq  web-frontend  2     1        1000   146.20    0.15

```
### Odditites and sorting
given the listed output above the optional --oddities flag picks out the containers that have a high cpu usage when compared to the other containers listed we also sort the list in descending order by the %REQ column
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend --oddities --sort '!%REQ'
T  PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
S  demo-odd-cpu-5f947f9db4-fg9wj  web-frontend  107   1        1000   10640.05  10.64
S  demo-odd-cpu-5f947f9db4-tthwv  web-frontend  105   1        1000   10426.41  10.43
S  demo-odd-cpu-5f947f9db4-dldv5  web-frontend  90    1        1000   8947.00   8.95

```
### Pod volumes
list all container volumes with mount points
``` shell
$ kubectl-ice volumes web-pod
CONTAINER    VOLUME                 TYPE       BACKING           SIZE  RO    MOUNT-POINT
app-init     kube-api-access-8z4tk  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-watcher  kube-api-access-8z4tk  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-broken   kube-api-access-8z4tk  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
myapp        app                    ConfigMap  app.py            -     false /myapp/
myapp        kube-api-access-8z4tk  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount

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
T  PODNAME                           CONTAINER       USED    REQUEST  LIMIT  %REQ    %LIMIT
S  demo-random-cpu-669b7888b9-46tzm  nginx           2.78Mi  1M       256M   291.23  1.14
S  demo-random-cpu-669b7888b9-56hbd  nginx           2.77Mi  1M       256M   290.00  1.13
S  demo-random-cpu-669b7888b9-jl4vg  nginx           2.82Mi  1M       256M   295.32  1.15
S  demo-random-cpu-669b7888b9-rbs6d  nginx           2.80Mi  1M       256M   294.09  1.15

```
