#!/bin/bash

## This script is not mantatory, it is here just to explain how we created 
## the Deployement resource mutator for Memcached
## $ ./gen-res.sh Deployment MCDeployment

type="$1" # Resource type (Secret, Service, ConfigMap, StatefulSet, ...)
tplFileOrName="$2"  # Either a YAML Template file of resource's initial data at deployment time or just the resource name.

#TODO: Pick up these parameters values from PROJECT file at the project's root folder
group="cache"
kind="Memcached"
version="v1alpha1"
path="github.com/example/memcached-operator"

# Generates an OKT resource from a template in YAML
okt-gen-resource -group ${group} -kind ${kind}  -path ${path} -type ${type} -version ${version} ${tplFileOrName}

exit $?

