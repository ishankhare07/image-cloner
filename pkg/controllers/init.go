package controllers

import (
	config "github.com/ishankhare07/image-cloner/pkg/config"
	"github.com/ishankhare07/image-cloner/pkg/registry"
)

var PrivateRegistryClient registry.Client
var UpstreamRegistryPool registry.ClientPool

func init() {
	PrivateRegistryClient = registry.NewRegistryClient(
		registry.WithName(config.PrivateRepoConfig.URI),
		registry.WithAuth(config.PrivateRepoConfig.ToAuthConfig()))

	UpstreamRegistryPool = make(registry.ClientPool)
}
