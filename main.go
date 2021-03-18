package main

import (
	"os"

	"github.com/ishankhare07/image-cloner/pkg/controllers"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func init() {
	log.SetLogger(zap.New())
}

func main() {
	entryLog := log.Log.WithName("entrypoint")

	entryLog.Info("setting up manager")

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to setup the controller manager")
		os.Exit(1)
	}

	if err = (&controllers.DeploymentReconciler{}).RegisterWithManager(mgr); err != nil {
		entryLog.Error(err, "unable to create controller", "controllers", "deployment")
		os.Exit(1)
	}

	entryLog.Info("Starting manager")
	if err = mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "problem starting manager")
		os.Exit(1)
	}
}
