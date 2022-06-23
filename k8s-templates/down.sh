#!/bin/bash

kubectl delete -f ./demo-probe.yml
kubectl delete -f ./demo-random-cpu.yml
kubectl delete -f ./demo-odd-cpu.yml
kubectl delete -f ./demo-memory.yml
kubectl delete -f ./demo-pod.yml
kubectl delete -f ./configmap.yml
kubectl delete -f ./namespace.yml

kubectl delete -n single-pods -f ./demo-probe.yml
kubectl delete -n cpu-demo -f ./demo-random-cpu.yml
kubectl delete -n resource-demo -f ./demo-odd-cpu.yml
kubectl delete -n resource-demo -f ./demo-memory.yml
kubectl delete -n single-pods -f ./demo-pod.yml