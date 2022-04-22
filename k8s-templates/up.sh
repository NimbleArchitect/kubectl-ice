#!/bin/bash

kubectl apply -f ./namespace.yml
kubectl apply -f ./configmap.yml
kubectl apply -f ./demo-pod.yml
kubectl apply -f ./demo-memory.yml
kubectl apply -f ./demo-odd-cpu.yml
kubectl apply -f ./demo-random-cpu.yml
kubectl apply -f ./demo-probe.yml
