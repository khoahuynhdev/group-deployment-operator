/*
Copyright 2024.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	multideploymentv1 "xplatform.api/group_deployment/api/v1"
)

// GroupDeploymentReconciler reconciles a GroupDeployment object
type GroupDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=multideployment.xplatform.api,resources=groupdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=multideployment.xplatform.api,resources=groupdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=multideployment.xplatform.api,resources=groupdeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GroupDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *GroupDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	groupDeployment := &multideploymentv1.GroupDeployment{}
	if err := r.Get(ctx, req.NamespacedName, groupDeployment); err != nil {
		log.Error(err, "Failed to get object")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	} else {
		if err = r.EnsureDeployment(ctx, groupDeployment); err != nil {
			log.Error(err, "Failed to ensure deployment")
			return ctrl.Result{}, err
		}
	}
	// create deployments
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GroupDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&multideploymentv1.GroupDeployment{}).
		Complete(r)
}
