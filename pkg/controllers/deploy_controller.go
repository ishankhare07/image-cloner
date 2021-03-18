package controllers

import (
	"context"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DeploymentReconciler struct {
	client client.Client
	log    logr.Logger
	scheme *runtime.Scheme
}

// Reconciler performs a full reconciliation for the object referred to by the Request.
// The Controller will requeue the Request to be processed again if an error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	reqLogger := r.log.WithValues("deploy", req.NamespacedName)

	if req.Namespace == "kube-system" {
		// ignore
		return reconcile.Result{}, nil
	}

	reqLogger.Info("event received for deploy", "info", req.NamespacedName)

	return ctrl.Result{}, nil
}

func (r *DeploymentReconciler) RegisterWithManager(mgr ctrl.Manager) error {
	r.client = mgr.GetClient()
	r.log = ctrl.Log.WithName("controller").WithName("deployment")
	r.scheme = mgr.GetScheme()

	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}
