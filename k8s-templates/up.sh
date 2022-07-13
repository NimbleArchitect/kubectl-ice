#!/bin/bash

if [[ "$1" == "ice" ]]; then
    kubectl apply -f ./namespace.yml
    kubectl apply -n ice -f ./configmap.yml
    
    kubectl apply -n ice -f ./demo-pod.yml
    kubectl apply -n ice -f ./demo-memory.yml
    kubectl apply -n ice -f ./demo-odd-cpu.yml
    kubectl apply -n ice -f ./demo-random-cpu.yml
    kubectl apply -n ice -f ./demo-probe.yml
else
    kubectl apply -f ./namespace.yml
    kubectl apply -n cpu-demo -f ./configmap.yml
    kubectl apply -n single-pods -f ./configmap.yml
    kubectl apply -n resource-demo -f ./configmap.yml

    kubectl apply -n single-pods -f ./demo-pod.yml
    kubectl apply -n resource-demo -f ./demo-memory.yml
    kubectl apply -n resource-demo -f ./demo-odd-cpu.yml
    kubectl apply -n cpu-demo -f ./demo-random-cpu.yml
    kubectl apply -n single-pods -f ./demo-probe.yml
fi
