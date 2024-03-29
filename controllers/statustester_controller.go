/*

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

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	azurev1 "github.com/Azure/statustest/api/v1"
	v1 "github.com/Azure/statustest/api/v1"
)

// StatusTesterReconciler reconciles a StatusTester object
type StatusTesterReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=azure.microsoft.com,resources=statustesters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=azure.microsoft.com,resources=statustesters/status,verbs=get;update;patch

func (r *StatusTesterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("statustester", req.NamespacedName)

	instance := &v1.StatusTester{}

	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		log.Info("couldnt get the resource")
		return ctrl.Result{}, err
	}

	instance.Status.Provisioned = true

	return ctrl.Result{}, r.Status().Update(ctx, instance)
}

func (r *StatusTesterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&azurev1.StatusTester{}).
		Complete(r)
}
