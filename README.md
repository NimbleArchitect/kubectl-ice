# kubectl-pod

This plugin shows useful information about the containers inside a pod

## Usage

Use `kubectl pod help` for help
```
kubectl pod cpu        # return cpu requests size and limits of each container
kubectl pod help       # Help about any command
kubectl pod image      # list the image name and pull status for each container
kubectl pod ip         # list ip addresses of all pods in the namespace listed
kubectl pod memory     # return memory requests size and limits of each container
kubectl pod restarts   # show restart counts for each container in a named pod
kubectl pod stats      # list resource usage of each container in a pod
kubectl pod status     # list status of each container in a pod
kubectl pod volumes    # list all container volumes with mount points
```

