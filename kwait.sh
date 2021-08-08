#!/bin/bash

echo "Wait fo ReconciliationSuccess condition..."
kubectl wait --for=condition=ReconciliationSuccess Memcached/memcached-sample

exit $?

