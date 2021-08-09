/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cachev1alpha1 "github.com/tapairmax/memcached-operator-with-okt/api/v1alpha1"

	//ADDED4OKT
	oktreconciler "gitlab.tech.orange/dbmsprivate/operators/okt/reconciler"
	oktengines "gitlab.tech.orange/dbmsprivate/operators/okt/reconciler/engines"
	okterr "gitlab.tech.orange/dbmsprivate/operators/okt/results"
)

// ADDED4OKT
var parameters map[string]string = map[string]string{
	"ContainerPort": "11211",
}

// MemcachedReconciler reconciles a Memcached object
/* REPLACEMENT4OKT
type MemcachedReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}
*/
// ADDED4OKT
type MemcachedReconciler struct {
	oktreconciler.AdvancedObject // Reconciler type

	// My resources
	CR         cachev1alpha1.Memcached
	deployment *ResourceMCDeploymentMutator
}

// Blank assignement to check type
var _ oktengines.StepperEngineHook = &MemcachedReconciler{}

//+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

/* REPLACED4OKT: OKT Reconciler let you manage the Reconciliation differently

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Memcached object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("memcached", req.NamespacedName)

	// Fetch the Memcached instance
	memcached := &cachev1alpha1.Memcached{}
	err := r.Get(ctx, req.NamespacedName, memcached)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Memcached resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Memcached")
		return ctrl.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: memcached.Name, Namespace: memcached.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForMemcached(memcached)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Ensure the deployment size is the same as the spec
	size := memcached.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		// Ask to requeue after 1 minute in order to give enough time for the
		// pods be created on the cluster side and the operand be able
		// to do the next update step accurately.
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	// Update the Memcached status with the pod names
	// List the pods for this memcached's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(memcached.Namespace),
		client.MatchingLabels(labelsForMemcached(memcached.Name)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "Memcached.Namespace", memcached.Namespace, "Memcached.Name", memcached.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, memcached.Status.Nodes) {
		memcached.Status.Nodes = podNames
		err := r.Status().Update(ctx, memcached)
		if err != nil {
			log.Error(err, "Failed to update Memcached status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}
*/
// EnterInState is the callback used by the Stepper engine when the Reconciliation enters in a new state
func (r *MemcachedReconciler) EnterInState(engine *oktengines.Stepper) {
	var err error

	switch engine.GetState() {
	// Main course states
	case "CRChecker":
		// We are here because the CR is picked from Cluster and any validation
		// webhooks, when they exists, are passed successfully
		// Perform here additional controls if any ...
	case "ObjectsGetter":
		// Creates the resource mutator object for the Memcached "Deployment"
		if r.deployment, err = NewResourceMCDeploymentMutator(&r.CR, r.Client, r.CR.GetNamespace(), r.CR.Name); err != nil {
			// Add error to the reconciliation results
			r.Results.AddOp(r.deployment, okterr.OperationResultRegistrationAborted, err, 0)
			return
		}
		r.RegisterResource(r.deployment)
		// Continue, hereafter, to register any other resource mutator the same way....

	case "Mutator":
		// Mutate all registered resources by calling their Mutation methods
		r.MutateAllResources(false)
	case "Updater":
		r.CreateOrUpdateAllResources(0, false) // Only for modified resources

		// Update the CR Status with the Node list (as in the original Memcached sample implementation)
		if r.CR.Status.Nodes, err = r.getPodNamesList(); err != nil {
			r.Results.AddOp(nil, okterr.OperationResultCRUDError, err, 0)
			return
		}
		//r.Results.DisplayOpList(r.Log)
	case "SuccessManager":
		r.ManageSuccess() // Will save the CR Status as well

	// Debranching states (exiting from the main course)
	case "CRFinalizer": // To come here, the CR must have a finalizer (in Metadata) with the same name as himself
		// Do here your custom finalization stuff
		r.Results.AddOpSuccess(nil, "Successfully finalized memcached")
	case "ErrorManager":
		// Take care of the CR Status and throttle repeated errors (if any) by
		// requeueing with a growing duration
		r.ManageError()
	default:
	}
}

/* REPLACED4OKT
** This dode moved into DeployementMC.go file for this resource creation and mutation
// deploymentForMemcached returns a memcached Deployment object
func (r *MemcachedReconciler) deploymentForMemcached(m *cachev1alpha1.Memcached) *appsv1.Deployment {
	ls := labelsForMemcached(m.Name)
	replicas := m.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   "memcached:1.4.36-alpine",
						Name:    "memcached",
						Command: []string{"memcached", "-m=64", "-o", "modern", "-v"},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 11211,
							Name:          "memcached",
						}},
					}},
				},
			},
		},
	}
	// Set Memcached instance as the owner and controller
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}
*/

// labelsForMemcached returns the labels for selecting the resources
// belonging to the given memcached CR name.
func labelsForMemcached(name string) map[string]string {
	return map[string]string{"app": "memcached", "memcached_cr": name}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// getPodsList
func getPodsList(c client.Client, namespace, name string) ([]string, error) {
	// List the pods for this memcached's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(namespace),
		client.MatchingLabels(labelsForMemcached(name)),
	}
	if err := c.List(context.TODO(), podList, listOpts...); err != nil {
		return nil, err
	}

	return getPodNames(podList.Items), nil
}

// getPodNamesList
func (r *MemcachedReconciler) getPodNamesList() ([]string, error) {
	return getPodsList(r.Client, r.CR.Namespace, r.CR.Name)
}

// SetupWithManager sets up the controller with the Manager.
func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// ADDED4OKT
	// If we no CR Status to manage, the following line is ok
	//r.Init("dev", &r.CR, nil)
	// We choose to manage a Status condition => so Init with the Conditions list
	r.Init("dev", &r.CR, &r.CR.Status.Conditions)

	engine := oktengines.NewStepper(r)
	r.SetEngine(engine)
	r.Params = parameters // Not mandatory but having common operator's parameters shared at one place is better

	// Same as the standard way
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Memcached{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
