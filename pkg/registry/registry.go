package registry

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
)

type Interface interface {
	GetImage(imgName string) (v1.Image, error)
	ImageExists(imageName string) (bool, name.Reference, error)
	Upload(img v1.Image, imageName string) error
}

type RegistryClientOption func(*registryClient)

type registryClient struct {
	registryName  string
	authenticator authn.Authenticator
}

func (r *registryClient) GetImage(imgName string) (v1.Image, error) {
	var registryName string
	if r.registryName == "" {
		registryName = imgName
	} else {
		registryName = r.registryName + "/" + imgName
	}

	ref, err := name.ParseReference(registryName)
	if err != nil {
		return nil, err
	}

	return remote.Image(ref, remote.WithAuth(r.authenticator))
}

func (r *registryClient) Upload(img v1.Image, imageName string) error {
	exists, ref, err := r.ImageExists(imageName)
	if exists {
		return nil
	}

	if err != nil {
		return err
	}

	return remote.Write(ref, img, remote.WithAuth(r.authenticator))
}

func (r *registryClient) ImageExists(imageName string) (bool, name.Reference, error) {
	ref, err := name.ParseReference(r.registryName + "/" + imageName)
	if err != nil {
		return false, nil, err
	}

	_, err = remote.Image(ref, remote.WithAuth(r.authenticator))
	if err != nil {
		terr, _ := err.(*transport.Error)
		for _, ec := range terr.Errors {
			if ec.Code == transport.ManifestUnknownErrorCode {
				return false, nil, nil
			}
		}
	}

	return true, ref, nil
}

func WithName(registryName string) RegistryClientOption {
	return func(r *registryClient) {
		r.registryName = registryName
	}
}

func WithAuth(cfg authn.AuthConfig) RegistryClientOption {
	return func(r *registryClient) {
		r.authenticator = authn.FromConfig(cfg)
	}
}

func NewRegistryClient(registryClientOptions ...RegistryClientOption) Interface {
	const (
		defaultName = ""
	)

	r := &registryClient{
		registryName:  defaultName,
		authenticator: authn.Anonymous,
	}

	for _, opt := range registryClientOptions {
		opt(r)
	}

	return r
}
