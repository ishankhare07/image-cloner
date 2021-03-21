package controllers

import (
	"context"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

	instance := &appsv1.Deployment{}
	err := r.client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// deployment might have been delete by now
			return ctrl.Result{}, nil
		}
		reqLogger.Error(err, "unable to get the deployment instance")
		return ctrl.Result{}, err
	}

	reqLogger.Info("deploy status details", "status", instance.Status)

	reqLogger.Info("images in deploy")

	var initContainersUpdated, containersUpdated bool
	initContainersUpdated, result, err := Cloner(reqLogger, instance.Spec.Template.Spec.InitContainers)
	if err != nil {
		return result, err
	}

	containersUpdated, result, err = Cloner(reqLogger, instance.Spec.Template.Spec.Containers)
	if err != nil {
		return result, err
	}

	if initContainersUpdated || containersUpdated {
		// a change was made in the deploy, hence update the object
		err := r.client.Update(ctx, instance)
		if err != nil {
			reqLogger.Error(err, "cannot update deploy")
			return ctrl.Result{}, err
		}
	}

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
