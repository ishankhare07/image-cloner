package utils

import (
	"fmt"
	"strings"
)

const (
	DefaultRegistry    = 1
	OrgRepo            = 2
	ThirdPartyRegistry = 3
)

func ExtractImageName(repoName string) string {
	splits := strings.Split(repoName, "/")

	return splits[len(splits)-1]
}

type RegistryInfo struct {
	registryName string
	repoName     string
	imageName    string
}

// GetNameForClient returns the registry name for clients
// along with the repo information if any. This can be
// used for initializing the RegistryClients' WithName
// options to set the right registry references
func (r *RegistryInfo) GetNameForClient() string {
	if r.repoName == "" {
		// case of DefaultRegistry
		return r.registryName
	}

	// case of OrgRepo or ThirdPartyRegistry
	return fmt.Sprintf("%s/%s", r.registryName, r.repoName)
}

func (r *RegistryInfo) GetImageName() string {
	return r.imageName
}

// ExtractRegistryInfo is supposed to take in the raw image
// string specified in the pod templates and returns a
// convenience object that can be further used to initialize
// respective RegistryClients and PrivateRegistryClients
func ExtractRegistryInfo(repoName string) *RegistryInfo {
	registry := "index.docker.io"
	repo := ""
	image := ""

	splits := strings.Split(repoName, "/")
	switch len(splits) {
	case DefaultRegistry:
		// is image uses a top level registry on docker hub
		// not typically hosted under a particular repo
		image = repoName
	case OrgRepo:
		// this image uses a registry on docker hub but uses
		// their own repo to host the image
		repo, image = splits[0], splits[1]
	case ThirdPartyRegistry:
		// this image uses a 3rd party registry
		// these are assumed to be organized
		// under their respective repos
		registry, repo, image = splits[0], splits[1], splits[2]
	}

	return &RegistryInfo{
		registryName: registry,
		repoName:     repo,
		imageName:    image,
	}
}
