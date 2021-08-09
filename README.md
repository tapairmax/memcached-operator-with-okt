## Intro

+ Note that we de-activated the webhooks in the `main.go` file. Uncomment to re-activate them if you have generated the expected certificates.
+ The `go.mod` file is updated to depend on OKT. A replace rule is added to use OKT locally if needed. The line is commented in order to pick up [OKT on Orange's Gitlab](https://gitlab.....orange/dbmsprivate/operators/okt) directly instead.
+ Have a look of the integration of the OKT Reconciler in `controller/memcached_controller.go` to perform the reconciliation
+ The Deployment resource files (`controller/MCDeployment.go` and `controller/DeploymentStub.go`) have been created thanks to the `gen-res.sh` shell script in the Controller folder. 
+ `controller/MCDeployment.go` is the file where you implement your own mutation for this resource throug 3 dedicated GO methods to customize:
  + 1-define the scope of data to put in the hash computation for change detection, 
  + 2-fill initial data in the 'Expected' GO struct object (to use to update the K8S cluster),  
  + 3-apply CR values in the 'Expected' object. 
+ The `controller/DeploymentStub.go` is a wrapper on a K8S Deployment resource/version. It would be shared among several Deployment resources used by this controller, if any. This file is not intended to be modified. 
+ The `gen-res.sh` shell depends on the installation of the `okt-gen-res` binary in your PATH. To install it, download OKT repo and install the command thanks to the Makefile. It will install the command in $HOME/bin (see OKT documentation).


## MemcaOperatorSDK Memcached project Init

    $ make generate
    $ make manifests

## MemcaOperatorSDK Memcached project Local Run

    $ make install run

## MemcaOperatorSDK Memcached project Wait for Reconciliation success

    $ kubectl wait --for=condition=ReconciliationSuccess Memcached/memcached-sample

## MemcaOperatorSDK Memcached project Cleanup

    $ kubectl delete -f config/samples/cache_v1alpha1_memcached.yaml
    $ make undeploy


## Application Lifecyle concerns thanks to a State Machine

Please see the markdown [here](https://github.com/tapairmax/memcached-operator-with-okt/blob/master/controllers/LATER-AppStateMachine.md)
