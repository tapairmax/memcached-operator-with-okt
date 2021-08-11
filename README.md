## Intro

This Operator is only usable with the Operator Karma Tools (OKT) which is still not available.

So it'll never work as is without OKT. This code just to illustrate an implementation to compare with the original Memcached Operator implemented with the OperatorSDK alone.

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


## Usage

### Clone OKT and Memcached-with-OKT repositries

        $ cd $HOME/go/src
	$ git clone https://github.com/Orange-OpenSource/Operators-Karma-Tools.git
	$ git clone https://github.com/tapairmax/memcached-operator-with-okt.git

### Memcached project Init

    $ make generate
    $ make manifests

### Memcached project Local Run and example of output

    $ make install run
    ....
    2021-...Z        INFO    controllers.Memcached.ENV=dev   History: >CRChecker>ObjectsGetter>Mutator>Updater>SuccessManager>End
    2021-...Z        INFO    controllers.Memcached.ENV=dev   Op: CR is succesfully picked up on Cluster      {"res": "Memcached/memcached-sample"}
    2021-...Z        INFO    controllers.Memcached.ENV=dev   Op: resource registration success       {"res": "Deployment/memcached-sample"}
    2021-...Z        INFO    controllers.Memcached.ENV=dev   Op: resource mutation success   {"res": "Deployment/memcached-sample"}
    2021-...Z        INFO    controllers.Memcached.ENV=dev   Op: resource unchanged  {"res": "Deployment/memcached-sample"}
    2021-...Z        INFO    controllers.Memcached.ENV=dev   Op: status updated      {"res": "Memcached/memcached-sample"}
    2021-...Z        INFO    controllers.Memcached.ENV=dev   Consolidated requeue duration: 0 seconds
    2021-...Z        INFO    controllers.Memcached.ENV=dev   1 CR is succesfully picked up on Cluster
    2021-...Z        INFO    controllers.Memcached.ENV=dev   1 resource registration success
    2021-...Z        INFO    controllers.Memcached.ENV=dev   1 resource mutation success
    2021-...Z        INFO    controllers.Memcached.ENV=dev   1 resource unchanged
    2021-...Z        INFO    controllers.Memcached.ENV=dev   1 status updated
    2021-...Z        INFO    controllers.Memcached.ENV=dev   5 Ops(s)
    2021-...Z        INFO    controllers.Memcached.ENV=dev   0 Error(s)


### Memcached CR creation

    $ vi config/samples/cache_v1alpha1_memcached.yaml
    $ kubectl create -f config/samples/cache_v1alpha1_memcached.yaml

### Memcached project Wait for Reconciliation success (after failure for example)

    $ kubectl wait --for=condition=ReconciliationSuccess Memcached/memcached-sample

### Memcached project Cleanup

    $ kubectl delete -f config/samples/cache_v1alpha1_memcached.yaml
    $ make undeploy


## Memcached Application Lifecyle concerns thanks to a State Machine

Please see the markdown [here](https://github.com/tapairmax/memcached-operator-with-okt/blob/master/controllers/LATER-AppStateMachine.md)
