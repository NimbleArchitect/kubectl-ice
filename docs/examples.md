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
demo-probe-76b66d5766-bft56  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-bft56  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-bft56  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/
demo-probe-76b66d5766-qx5f2  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-qx5f2  web-frontend  readiness 5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-qx5f2  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/

```
### Named containers
the optional container flag (-c) searchs all selected pods and lists only containers that match the name web-frontend
``` shell
$ kubectl-ice command -c web-frontend
T  PODNAME                           CONTAINER     COMMAND                                       ARGUMENTS
S  demo-memory-7ddb58cd5b-5psfr      web-frontend  python /myapp/halfmemapp.py                   -
S  demo-memory-7ddb58cd5b-6sjtl      web-frontend  python /myapp/halfmemapp.py                   -
S  demo-memory-7ddb58cd5b-fvhl7      web-frontend  python /myapp/halfmemapp.py                   -
S  demo-memory-7ddb58cd5b-zn5qg      web-frontend  python /myapp/halfmemapp.py                   -
S  demo-odd-cpu-5f947f9db4-72f89     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-odd-cpu-5f947f9db4-78vks     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-odd-cpu-5f947f9db4-85zq2     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-odd-cpu-5f947f9db4-872gg     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-odd-cpu-5f947f9db4-dwmjd     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-odd-cpu-5f947f9db4-g8468     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-odd-cpu-5f947f9db4-gjcc2     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-odd-cpu-5f947f9db4-hg9s9     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-odd-cpu-5f947f9db4-kjx4r     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-odd-cpu-5f947f9db4-klsp4     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-odd-cpu-5f947f9db4-n2gv7     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-odd-cpu-5f947f9db4-phnjp     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-odd-cpu-5f947f9db4-qt4sg     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-odd-cpu-5f947f9db4-tds8q     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-odd-cpu-5f947f9db4-trff6     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-odd-cpu-5f947f9db4-x5cm5     web-frontend  python /myapp/oddcpuapp.py                    -
S  demo-probe-76b66d5766-bft56       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0   -
S  demo-probe-76b66d5766-qx5f2       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0   -
S  demo-random-cpu-669b7888b9-74gt7  web-frontend  python /myapp/randomcpuapp.py                 -
S  demo-random-cpu-669b7888b9-9lcrs  web-frontend  python /myapp/randomcpuapp.py                 -
S  demo-random-cpu-669b7888b9-bs6rs  web-frontend  python /myapp/randomcpuapp.py                 -
S  demo-random-cpu-669b7888b9-mgvp7  web-frontend  python /myapp/randomcpuapp.py                 -

```
### Labels and containers
you can also search specific pods and list all containers with a specific name, in this example all pods with the label app=userandomcpu are searched and only the containers that match the name web-fronteend are shown
``` shell
$ kubectl-ice cpu -l app=userandomcpu -c web-frontend
PODNAME                           CONTAINER     USED  REQUEST  LIMIT  %REQ    %LIMIT
demo-random-cpu-669b7888b9-74gt7  web-frontend  122m  125m     1000m  96.97   12.12
demo-random-cpu-669b7888b9-9lcrs  web-frontend  129m  125m     1000m  102.65  12.83
demo-random-cpu-669b7888b9-bs6rs  web-frontend  149m  125m     1000m  119.12  14.89
demo-random-cpu-669b7888b9-mgvp7  web-frontend  208m  125m     1000m  165.89  20.74

```
### Container status
most commands work the same way including the status command which also lets you see which container(s) are causing the restarts and by using the optional --previous flag you can view the containers previous exit code
``` shell
$ kubectl-ice status -l app=myapp --previous
T  PODNAME  CONTAINER    STATE       REASON              EXIT-CODE  SIGNAL  TIMESTAMP                      MESSAGE
S  web-pod  app-broken   Terminated  Error               1          0       2022-05-31 19:08:58 +0100 BST  -
S  web-pod  app-watcher  Terminated  Error               2          0       2022-05-31 19:10:04 +0100 BST  -
S  web-pod  myapp        Terminated  ContainerCannotRun  127        0       2022-05-31 19:08:59 +0100 BST  OCI runtime create failed: container_linux.go:380: starting container process caused: exec: "python /myapp/mainapp.py\nwith\nmultiline\nargument\n": stat python /myapp/mainapp.py
with
multiline
argument
: no such file or directory: unknown
I  web-pod  app-init     -           -                   -          -       -                              -

```
### Container images
need to chack on the currently configured image versions use the image command
``` shell
$ kubectl-ice image -l app=userandomcpu
T  PODNAME                           CONTAINER       PULL          IMAGE
S  demo-random-cpu-669b7888b9-74gt7  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-74gt7  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-74gt7  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-9lcrs  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-9lcrs  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-9lcrs  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-bs6rs  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-bs6rs  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-bs6rs  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-mgvp7  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-mgvp7  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-mgvp7  init-myservice  IfNotPresent  busybox:1.28

```
### Advanced labels
return memory requests size and limits of each container where the pods have an app label that matches useoddcpu and the container name is equal to nginx
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-72f89  web-frontend  2m    1m       1000m  153.80    0.15
demo-odd-cpu-5f947f9db4-78vks  web-frontend  3m    1m       1000m  200.42    0.20
demo-odd-cpu-5f947f9db4-85zq2  web-frontend  2m    1m       1000m  162.59    0.16
demo-odd-cpu-5f947f9db4-872gg  web-frontend  106m  1m       1000m  10534.03  10.53
demo-odd-cpu-5f947f9db4-dwmjd  web-frontend  2m    1m       1000m  153.30    0.15
demo-odd-cpu-5f947f9db4-g8468  web-frontend  2m    1m       1000m  155.68    0.16
demo-odd-cpu-5f947f9db4-gjcc2  web-frontend  2m    1m       1000m  151.98    0.15
demo-odd-cpu-5f947f9db4-hg9s9  web-frontend  3m    1m       1000m  232.32    0.23
demo-odd-cpu-5f947f9db4-kjx4r  web-frontend  2m    1m       1000m  148.07    0.15
demo-odd-cpu-5f947f9db4-klsp4  web-frontend  94m   1m       1000m  9382.23   9.38
demo-odd-cpu-5f947f9db4-n2gv7  web-frontend  2m    1m       1000m  156.18    0.16
demo-odd-cpu-5f947f9db4-phnjp  web-frontend  2m    1m       1000m  148.46    0.15
demo-odd-cpu-5f947f9db4-qt4sg  web-frontend  3m    1m       1000m  209.04    0.21
demo-odd-cpu-5f947f9db4-tds8q  web-frontend  3m    1m       1000m  235.68    0.24
demo-odd-cpu-5f947f9db4-trff6  web-frontend  2m    1m       1000m  150.69    0.15
demo-odd-cpu-5f947f9db4-x5cm5  web-frontend  2m    1m       1000m  155.55    0.16

```
### Odditites and sorting
given the listed output above the optional --oddities flag picks out the containers that have a high cpu usage when compared to the other containers listed we also sort the list in descending order by the %REQ column
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend --oddities --sort '!%REQ'
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-872gg  web-frontend  106m  1m       1000m  10534.03  10.53
demo-odd-cpu-5f947f9db4-klsp4  web-frontend  94m   1m       1000m  9382.23   9.38

```
### Pod volumes
list all container volumes with mount points
``` shell
$ kubectl-ice volumes web-pod
CONTAINER    VOLUME                 TYPE       BACKING           SIZE  RO    MOUNT-POINT
app-init     kube-api-access-jg4mh  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-watcher  kube-api-access-jg4mh  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-broken   kube-api-access-jg4mh  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
myapp        app                    ConfigMap  app.py            -     false /myapp/
myapp        kube-api-access-jg4mh  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount

```
### Pod exec command
retrieves the command line and any arguments specified at the container level
``` shell
$ kubectl-ice command web-pod
T  CONTAINER    COMMAND                                             ARGUMENTS
S  app-watcher  python /myapp/mainapp.py                            -
S  app-broken   sh -c sleep 2; exit 1                               -
S  myapp        python /myapp/mainapp.py with multiline argument    -
I  app-init     sh -c sleep 2; exit 0                               -

```
### Excluding rows
use the --match flag to show only the output rows where the used memory column is greater than or equal to 1MB, this has the effect of exclusing any row where the used memory column is currently under 1MB, the value 1 can be replace with any whole number in megabytes, to show only used memory greater than 1GB you would replace 1 with 1000
``` shell
$ kubectl-ice mem -l app=userandomcpu --match 'used>=1'
PODNAME                           CONTAINER     USED    REQUEST  LIMIT  %REQ    %LIMIT
demo-random-cpu-669b7888b9-74gt7  web-frontend  5.36Mi  1M       256M   561.56  2.19
demo-random-cpu-669b7888b9-74gt7  nginx         3.73Mi  1M       256M   391.58  1.53
demo-random-cpu-669b7888b9-9lcrs  web-frontend  7.33Mi  1M       256M   768.41  3.00
demo-random-cpu-669b7888b9-9lcrs  nginx         3.58Mi  1M       256M   375.19  1.47
demo-random-cpu-669b7888b9-bs6rs  web-frontend  7.93Mi  1M       256M   831.90  3.25
demo-random-cpu-669b7888b9-bs6rs  nginx         3.79Mi  1M       256M   396.90  1.55
demo-random-cpu-669b7888b9-mgvp7  web-frontend  8.44Mi  1M       256M   884.74  3.46
demo-random-cpu-669b7888b9-mgvp7  nginx         3.56Mi  1M       256M   373.56  1.46

```
