#!/bin/bash

## This script is not mantatory, it is here just to explain how we created 
## the Deployement resource mutator for Memcached
## $ ./gen-res.sh <path_to>/PROJECT Deployment MCDeployment
## $ ./gen-res.sh <path_to>/PROJECT --app MCAppStateMachine


projFile="$1"
[ ! -f "$1" ] && echo "Error: the PROJECT file ${projFile} can not be read" && exit 1
shift

type="$1" # Resource type (Secret, Service, ConfigMap, StatefulSet, Ingress, Pod, ...). List with "okt-gen-resource --list"
tplFileOrName="$2"  # Either a YAML Template file of resource's initial data at deployment time or just the resource name.


trim() {
    local trimmed="$1"
    # Strip leading space.
    trimmed="${trimmed## }"
    # Strip trailing space.
    trimmed="${trimmed%% }"

    echo -n "$trimmed"
}

getProjectVal() {
    local val="$(grep "$1" ${projFile} | awk -F: '{print $2}')"
    echo -n "$(trim "${val}")"
}

# Pick up these parameters values from PROJECT file at the project's root folder
group="$(getProjectVal group)"
kind="$(getProjectVal kind)"
version="$(getProjectVal "  version")"
path="$(getProjectVal path)"

# Generates an OKT resource from a template in YAML
echo okt-gen-resource -group ${group} -kind ${kind}  -path ${path} -type ${type} -version ${version} ${tplFileOrName}
okt-gen-resource -group ${group} -kind ${kind}  -path ${path} -type ${type} -version ${version} ${tplFileOrName}

exit $?

