package registry

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type Interface interface {
	GetImage(imgName string) (v1.Image, error)
	Upload(img v1.Image, imageName string) error
}

type RegistryOption func(*registry)

type registry struct {
	registry      string
	Reference     name.Reference
	Authenticator authn.Authenticator
}

func (r *registry) GetImage(imgName string) (v1.Image, error) {
	var registryName string
	if r.registry == "" {
		registryName = imgName
	} else {
		registryName = r.registry + "/" + imgName
	}

	ref, err := name.ParseReference(registryName)
	if err != nil {
		return nil, err
	}

	return remote.Image(ref, remote.WithAuth(r.Authenticator))
}

func (r *registry) Upload(img v1.Image, imageName string) error {
	ref, err := name.ParseReference(r.registry + "/" + imageName)
	if err != nil {
		return err
	}

	return remote.Write(ref, img, remote.WithAuth(r.Authenticator))
}

func WithName(registryName string) RegistryOption {
	return func(r *registry) {
		r.registry = registryName
	}
}

func WithAuth(cfg authn.AuthConfig) RegistryOption {
	return func(r *registry) {
		r.Authenticator = authn.FromConfig(cfg)
	}
}

func NewRegistry(registryOptions ...RegistryOption) Interface {
	const (
		defaultName = ""
	)

	r := &registry{
		registry:      defaultName,
		Authenticator: authn.Anonymous,
	}

	for _, opt := range registryOptions {
		opt(r)
	}

	return r
}
