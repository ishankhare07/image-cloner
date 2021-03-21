package controllers

import (
	"github.com/go-logr/logr"
	"github.com/ishankhare07/image-cloner/pkg/registry"
	v1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Cloner(logger logr.Logger, containers []v1.Container) (bool, reconcile.Result, error) {
	// updated tracks if any of the images were uploaded for the given workload
	// or the image field was changed to an already uploaded image
	// in which case the workload object has to be reloaded
	updated := false
	for i, container := range containers {
		logger.Info("init container", "image", container.Image)

		registryInfo := registry.ExtractRegistryInfo(container.Image)

		logger.Info("********************** before comparision ************************************",
			"GetNameForClient", registryInfo.GetNameForClient(),
			"GetRegistryName", PrivateRegistryClient.GetRegistryName())
		if registryInfo.GetNameForClient() == PrivateRegistryClient.GetRegistryName() {
			logger.Info("image already cached in private repo skipping")
			continue
		}

		var client registry.Client
		var ok bool

		if client, ok = UpstreamRegistryPool[registryInfo.GetNameForClient()]; !ok {
			logger.Info("client does not exist", "repo", registryInfo.GetNameForClient())
			logger.Info("creating one now")
			client = registry.NewRegistryClient(
				registry.WithName(registryInfo.GetNameForClient()))
			UpstreamRegistryPool[registryInfo.GetNameForClient()] = client
			logger.Info("successfully created client", "repo", registryInfo.GetNameForClient())
		}

		img, err := client.GetImage(registryInfo.GetImageName())
		if err != nil {
			logger.Error(err, "unable to fetch image from upstream registry", registryInfo.GetNameForClient(), registryInfo.GetImageName())
			return updated, ctrl.Result{}, err
		}

		logger.Info("successfully fetched image", "name", registryInfo.GetImageName())

		logger.Info("trying to upload image to private repo")
		ok, err = PrivateRegistryClient.Upload(img, registryInfo.GetImageName())
		if err != nil {
			logger.Error(err, "unable to upload image to private registry")
			return updated, ctrl.Result{}, err
		}

		logger.Info("returned from upload function", "ok", ok)

		desiredImageName := PrivateRegistryClient.GetRegistryName() + "/" + registryInfo.GetImageName()
		if ok || container.Image != desiredImageName {
			logger.Info("successfully uploaded image", "image", registryInfo.GetImageName())
			logger.Info("updating container image", "container name", container.Name)
			updated = true
			containers[i].Image = desiredImageName
			logger.Info("update container image name", "new name", desiredImageName)
		}
	}

	return updated, ctrl.Result{}, nil
}
