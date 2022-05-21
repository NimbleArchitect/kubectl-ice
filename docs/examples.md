### Single pod info
Shows the currently used memory along with the configured memory requests and limits of all containers (side cars) in the pod named web-pod
``` shell
$ kubectl-ice memory web-pod
CONTAINER    USED    REQUEST  LIMIT  %REQ    %LIMIT
app-watcher  0       1M       512M   -       -
app-broken   0       1M       512M   -       -
myapp        0.01Gi  1M       256M   815.10  3.18

```
### Using labels
using labels you can search all pods that are part of a deployment where the label app matches demoprobe and list selected information about the containers in each pod, this example shows the currently configured probe information and gives details of configured startup, readiness and liveness probes of each container
``` shell
$ kubectl-ice probes -l app=demoprobe
PODNAME                      CONTAINER     PROBE     DELAY  PERIOD  TIMEOUT  SUCCESS  FAILURE  CHECK    ACTION
demo-probe-76b66d5766-j9wnm  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-j9wnm  web-frontend  liveness  5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-j9wnm  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/
demo-probe-76b66d5766-ksn5t  web-frontend  liveness  10     5       1        1        3        Exec     exit 0
demo-probe-76b66d5766-ksn5t  web-frontend  liveness  5      5       1        1        3        Exec     cat /tmp/health
demo-probe-76b66d5766-ksn5t  nginx         liveness  60     60      1        1        8        HTTPGet  http://:80/

```
### Named containers
the optional container flag (-c) searchs all selected pods and lists only containers that match the name web-frontend
``` shell
$ kubectl-ice command -c web-frontend
T  PODNAME                           CONTAINER     COMMAND                                      ARGUMENTS
S  demo-memory-7ddb58cd5b-5psfr      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-6sjtl      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-fvhl7      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-memory-7ddb58cd5b-zn5qg      web-frontend  python /myapp/halfmemapp.py                  -
S  demo-odd-cpu-5f947f9db4-459t8     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-6mlk9     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-7xcqw     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-8fc4c     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-9x5mb     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-bxchg     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-fsccd     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-gtlcl     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-j882g     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-mqwnd     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-qh7gk     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-rcxjq     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-rrj7c     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-rtxlm     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-xs2gs     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-odd-cpu-5f947f9db4-zx5c8     web-frontend  python /myapp/oddcpuapp.py                   -
S  demo-probe-76b66d5766-j9wnm       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
S  demo-probe-76b66d5766-ksn5t       web-frontend  sh -c touch /tmp/health; sleep 2000; exit 0  -
S  demo-random-cpu-669b7888b9-2gljr  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-hd4vv  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-lxl4b  web-frontend  python /myapp/randomcpuapp.py                -
S  demo-random-cpu-669b7888b9-wvrmq  web-frontend  python /myapp/randomcpuapp.py                -

```
### Labels and containers
you can also search specific pods and list all containers with a specific name, in this example all pods with the label app=userandomcpu are searched and only the containers that match the name web-fronteend are shown
``` shell
$ kubectl-ice cpu -l app=userandomcpu -c web-frontend
PODNAME                           CONTAINER     USED  REQUEST  LIMIT  %REQ    %LIMIT
demo-random-cpu-669b7888b9-2gljr  web-frontend  0m    125m     1000m  -       -
demo-random-cpu-669b7888b9-hd4vv  web-frontend  0m    125m     1000m  -       -
demo-random-cpu-669b7888b9-lxl4b  web-frontend  510m  125m     1000m  407.69  50.96
demo-random-cpu-669b7888b9-wvrmq  web-frontend  0m    125m     1000m  -       -

```
### Container status
most commands work the same way including the status command which also lets you see which container(s) are causing the restarts and by using the optional --previous flag you can view the containers previous exit code
``` shell
$ kubectl-ice status -l app=myapp --previous
T  PODNAME  CONTAINER    STATE       REASON  EXIT-CODE  SIGNAL  TIMESTAMP                      MESSAGE
S  web-pod  app-broken   Terminated  Error   1          0       2022-05-20 10:08:40 +0100 BST  -
S  web-pod  app-watcher  Terminated  Error   2          0       2022-05-20 10:04:54 +0100 BST  -
S  web-pod  myapp        -           -       -          -       -                              -
I  web-pod  app-init     -           -       -          -       -                              -

```
### Container images
need to chack on the currently configured image versions use the image command
``` shell
$ kubectl-ice image -l app=userandomcpu
T  PODNAME                           CONTAINER       PULL          IMAGE
S  demo-random-cpu-669b7888b9-2gljr  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-2gljr  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-2gljr  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-hd4vv  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-hd4vv  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-hd4vv  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-lxl4b  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-lxl4b  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-lxl4b  init-myservice  IfNotPresent  busybox:1.28
S  demo-random-cpu-669b7888b9-wvrmq  web-frontend    Always        python:latest
S  demo-random-cpu-669b7888b9-wvrmq  nginx           IfNotPresent  nginx:1.7.9
I  demo-random-cpu-669b7888b9-wvrmq  init-myservice  IfNotPresent  busybox:1.28

```
### Advanced labels
return memory requests size and limits of each container where the pods have an app label that matches useoddcpu and the container name is equal to nginx
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-459t8  web-frontend  2m    1m       1000m  168.67    0.17
demo-odd-cpu-5f947f9db4-6mlk9  web-frontend  110m  1m       1000m  10944.21  10.94
demo-odd-cpu-5f947f9db4-7xcqw  web-frontend  110m  1m       1000m  10991.84  10.99
demo-odd-cpu-5f947f9db4-8fc4c  web-frontend  2m    1m       1000m  161.52    0.16
demo-odd-cpu-5f947f9db4-9x5mb  web-frontend  2m    1m       1000m  163.67    0.16
demo-odd-cpu-5f947f9db4-bxchg  web-frontend  3m    1m       1000m  239.96    0.24
demo-odd-cpu-5f947f9db4-fsccd  web-frontend  115m  1m       1000m  11412.70  11.41
demo-odd-cpu-5f947f9db4-gtlcl  web-frontend  2m    1m       1000m  164.95    0.16
demo-odd-cpu-5f947f9db4-j882g  web-frontend  123m  1m       1000m  12257.90  12.26
demo-odd-cpu-5f947f9db4-mqwnd  web-frontend  3m    1m       1000m  258.07    0.26
demo-odd-cpu-5f947f9db4-qh7gk  web-frontend  2m    1m       1000m  176.79    0.18
demo-odd-cpu-5f947f9db4-rcxjq  web-frontend  3m    1m       1000m  273.59    0.27
demo-odd-cpu-5f947f9db4-rrj7c  web-frontend  121m  1m       1000m  12090.36  12.09
demo-odd-cpu-5f947f9db4-rtxlm  web-frontend  130m  1m       1000m  12917.46  12.92
demo-odd-cpu-5f947f9db4-xs2gs  web-frontend  2m    1m       1000m  164.32    0.16
demo-odd-cpu-5f947f9db4-zx5c8  web-frontend  2m    1m       1000m  175.37    0.18

```
### Odditites and sorting
given the listed output above the optional --oddities flag picks out the containers that have a high cpu usage when compared to the other containers listed we also sort the list in descending order by the %REQ column
``` shell
$ kubectl-ice cpu -l "app in (useoddcpu)" -c web-frontend --oddities --sort '!%REQ'
PODNAME                        CONTAINER     USED  REQUEST  LIMIT  %REQ      %LIMIT
demo-odd-cpu-5f947f9db4-rtxlm  web-frontend  130m  1m       1000m  12917.46  12.92
demo-odd-cpu-5f947f9db4-j882g  web-frontend  123m  1m       1000m  12257.90  12.26
demo-odd-cpu-5f947f9db4-rrj7c  web-frontend  121m  1m       1000m  12090.36  12.09
demo-odd-cpu-5f947f9db4-fsccd  web-frontend  115m  1m       1000m  11412.70  11.41
demo-odd-cpu-5f947f9db4-7xcqw  web-frontend  110m  1m       1000m  10991.84  10.99
demo-odd-cpu-5f947f9db4-6mlk9  web-frontend  110m  1m       1000m  10944.21  10.94
demo-odd-cpu-5f947f9db4-rcxjq  web-frontend  3m    1m       1000m  273.59    0.27
demo-odd-cpu-5f947f9db4-mqwnd  web-frontend  3m    1m       1000m  258.07    0.26
demo-odd-cpu-5f947f9db4-bxchg  web-frontend  3m    1m       1000m  239.96    0.24
demo-odd-cpu-5f947f9db4-qh7gk  web-frontend  2m    1m       1000m  176.79    0.18
demo-odd-cpu-5f947f9db4-zx5c8  web-frontend  2m    1m       1000m  175.37    0.18
demo-odd-cpu-5f947f9db4-459t8  web-frontend  2m    1m       1000m  168.67    0.17
demo-odd-cpu-5f947f9db4-gtlcl  web-frontend  2m    1m       1000m  164.95    0.16
demo-odd-cpu-5f947f9db4-xs2gs  web-frontend  2m    1m       1000m  164.32    0.16
demo-odd-cpu-5f947f9db4-9x5mb  web-frontend  2m    1m       1000m  163.67    0.16
demo-odd-cpu-5f947f9db4-8fc4c  web-frontend  2m    1m       1000m  161.52    0.16

```
### Pod volumes
list all container volumes with mount points
``` shell
$ kubectl-ice volumes web-pod
CONTAINER    VOLUME                 TYPE       BACKING           SIZE  RO    MOUNT-POINT
app-init     kube-api-access-4h6h2  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-watcher  kube-api-access-4h6h2  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
app-broken   kube-api-access-4h6h2  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount
myapp        app                    ConfigMap  app.py            -     false /myapp/
myapp        kube-api-access-4h6h2  Projected  kube-root-ca.crt  -     true  /var/run/secrets/kubernetes.io/serviceaccount

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
demo-random-cpu-669b7888b9-2gljr  nginx         4.57Mi  1M       256M   479.23  1.87
demo-random-cpu-669b7888b9-hd4vv  nginx         4.62Mi  1M       256M   484.56  1.89
demo-random-cpu-669b7888b9-lxl4b  web-frontend  7.05Mi  1M       256M   739.33  2.89
demo-random-cpu-669b7888b9-lxl4b  nginx         3.69Mi  1M       256M   386.66  1.51
demo-random-cpu-669b7888b9-wvrmq  nginx         4.60Mi  1M       256M   482.10  1.88

```
