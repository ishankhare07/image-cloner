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

type DaemonSetReconciler struct {
	client client.Client
	log    logr.Logger
	scheme *runtime.Scheme
}

// Reconciler performs a full reconciliation for the object referred to by the Request.
// The Controller will requeue the Request to be processed again if an error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (d *DaemonSetReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	reqLogger := d.log.WithValues("daemonset", req.NamespacedName)

	if req.Namespace == "kube-system" {
		// ignore
		return reconcile.Result{}, nil
	}

	reqLogger.Info("event received for daemonset", "info", req.NamespacedName)

	instance := &appsv1.DaemonSet{}
	err := d.client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// daemonset might have been deleted by now
			return ctrl.Result{}, err
		}
		reqLogger.Error(err, "unable to get the daemonset object")
		return ctrl.Result{}, err
	}

	reqLogger.Info("images in daemonset")
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
		// a change was made in the daemonset, hence update the object
		err := d.client.Update(ctx, instance)
		if err != nil {
			reqLogger.Error(err, "cannot update daemonset")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (d *DaemonSetReconciler) RegisterWithManager(mgr ctrl.Manager) error {
	d.client = mgr.GetClient()
	d.log = ctrl.Log.WithName("controller").WithName("daemonset")
	d.scheme = mgr.GetScheme()

	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.DaemonSet{}).
		Complete(d)
}
