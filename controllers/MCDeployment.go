package controllers

/**
*** THIS FILE IS AUTO-GENERATED BY THE OKT RESOURCE GENERATOR (CLI TOOL)
*** It implements the ResourceMutator object for a resource managed by the OKT reconciler.
*** It must be updated/customized with:
***   - The initial data that define this resource at deployement time (either a GO Struct or a YAML in string format)
***   - The custom of GetHashableRef() that will define the data in this resource on which the Hash computation will be done
***   - The 2 Mutation fonctions to fill both the expected resource object with the Initial data and with the CR values
**/

import (
	appapi "github.com/example/memcached-operator/api/v1alpha1"

	oktres "gitlab.tech.orange/dbmsprivate/operators/okt/resources"
	okthash "gitlab.tech.orange/dbmsprivate/operators/okt/tools/hash"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

/* Useless for this example, the resource is created thanks to a GO structure and not a YAML
func getTpl() string {
        yaml := `

`
        return yaml
}
*/

// 4OKT picked up from the Memcached example and moved here directly in the resource file
// deploymentForMemcached returns a memcached Deployment object
func deploymentForMemcached(m *appapi.Memcached) *appsv1.Deployment {
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
							ContainerPort: 11211, // TODO: Replace hard-coded value by r.GetData("ContainerPort") parameters
							Name:          "memcached",
						}},
					}},
				},
			},
		},
	}
	/* REMOVE4OKT (it is done automatically by OKT)
	   // Set Memcached instance as the owner and controller
	   ctrl.SetControllerReference(m, dep, r.Scheme)
	*/
	return dep
}


// ResourceMCDeploymentMutator xx
type ResourceMCDeploymentMutator struct {
	DeploymentResourceStub
	cr *appapi.Memcached
}

// blank assignment to verify this resource implements an OKT Resource
var _ oktres.Mutator = &ResourceMCDeploymentMutator{}

// NewResourceMCDeploymentMutator xx
func NewResourceMCDeploymentMutator(cr *appapi.Memcached, client k8sclient.Client, namespace, name string) (*ResourceMCDeploymentMutator, error) {
	res := &ResourceMCDeploymentMutator{cr: cr}

	if err := res.Init(client, namespace, name); err != nil {
		return nil, err
	}

	return res, nil
}

//--
// TODO: CUSTOMIZE HERE YOUR OWN MUTATIONS WITH DEFAULTS AND THE CUSTOM RESOURCE
//--

// GetHashableRef xx
// Note that a Spec reference can always be added by the helper. It is either the K8S Object's Spec or
// data as defined by OKT dictionary used by the resource code generator
func (r *ResourceMCDeploymentMutator) GetHashableRef() okthash.HashableRef {
	helper := r.GetHashableRefHelper()
	helper.AddMetaLabels()
	//helper.AddMetaLabelValues("app", "memcached_cr")
	//helper.AddMetaAnnotations()
	helper.AddUserData(&r.Expected.Spec)

	return helper
}

// MutateWithInitialData Initialize the Expected object with intial deployment data
func (r *ResourceMCDeploymentMutator) MutateWithInitialData() error {
	/* Useless here while we load initial data from a GO structure not a YAML string!
	   yaml := getTpl()
	   if err := r.CopyTpl(yaml, r.GetData()); err != nil {
	           return err
	   }
	*/
	// Alternatively, instead of a YAML data, use a GO structure
	dep := deploymentForMemcached(r.cr)
	// And copy (merge) it into the Expected resource.
	/***
	    dep.DeepCopyInto(&r.Expected)  DO NOT USE DEEP COPY!!! YOU'RE RECONCILIATION WILL LOOP FOR EVER
	**/
	if err := r.CopyGOStruct(dep); err != nil {
		return err
	}

	return nil
}

// MutateWithCR xx
func (r *ResourceMCDeploymentMutator) MutateWithCR() (requeueAfterSeconds uint16, err error) {
	// Apply CR values
	//r.Expected.Spec.xxx = r.cr.Spec.xxx

	size := r.cr.Spec.Size
	if *r.Expected.Spec.Replicas != size {
		/* Test to fall in error
		   if size == 4 {
		           return 10, errors.New("For some reason the Cluster can not scale to 4")
		   }
		*/
		r.Expected.Spec.Replicas = &size

		// Ask to requeue after 1 minute in order to give enough time for the
		// pods be created on the cluster side and the operand be able
		// to do the next update step accurately.
		return 60, nil
	}

	return 0, nil
}