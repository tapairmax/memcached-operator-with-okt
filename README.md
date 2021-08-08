## Intro

+ Note that we de-activated the webhooks in the `main.go` file. Uncomment to re-activate them if you have the expected certificates.
+ The `go.mod` file is updated to depend on OKT. A replace rule is added to use OKT locally if needed. The line is commented in order to pick up [OKT on Gitlab DIOD](https://gitlab.tech.orange/dbmsprivate/operators/okt) directly instead. 
+ The Deployment resource files (`controller/MCDeployment.go` and `controller/DeploymentStub.go`) have been created thanks to the `gen-res.sh` shell script in the Controller folder. 
+ `controller/MCDeployment.go` is the file where you implement your own mutation for this resource (1-define the scope of data to put in the hash computation for change detection, 2-fill initial data and 3-apply CR values). 
+ The `controller/DeploymentStub.go` is a wrapper on a K8S Deployment resource/version. It would be shared among several Deployment resources used by this controller, if any. This file is not intended to be modified. 
+ The `gen-res.sh` shell depends on the installation of the `okt-gen-res` binary in your PATH. To install it, download OKT repo and install the command thanks to the Makefile. It will install the command in $HOME/bin (see OKT documentation).


## Init

    $ make generate
    $ make manifests

## Local Run

    $ make install run

## Wait for Reconciliation success

    $ kubectl wait --for=condition=ReconciliationSuccess Memcached/memcached-sample

## Cleanup

    $ kubectl delete -f config/samples/cache_v1alpha1_memcached.yaml
    $ make undeploy

