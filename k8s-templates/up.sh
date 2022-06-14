#!/bin/bash

kubectl apply -f ./namespace.yml
kubectl apply -n cpu-demo -f ./configmap.yml
kubectl apply -n single-pods -f ./configmap.yml
kubectl apply -n resource-demo -f ./configmap.yml
# kubectl apply -n  -f ./configmap.yml
# kubectl apply -n cpu-demo -f ./configmap.yml

kubectl apply -n single-pods -f ./demo-pod.yml
kubectl apply -n resource-demo -f ./demo-memory.yml
kubectl apply -n resource-demo -f ./demo-odd-cpu.yml
kubectl apply -n cpu-demo -f ./demo-random-cpu.yml
kubectl apply -n single-pods -f ./demo-probe.yml
