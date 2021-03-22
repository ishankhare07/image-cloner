package repo

import (
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/spf13/viper"
)

type DockerRegistryConfig struct {
	URI      string
	Username string
	Password string
}

func (d *DockerRegistryConfig) ToAuthConfig() authn.AuthConfig {
	return authn.AuthConfig{
		Username: d.Username,
		Password: d.Password,
	}
}

var PrivateRepoConfig DockerRegistryConfig

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$PWD/config/")
	viper.AutomaticEnv()

	if val, ok := os.LookupEnv("DEV_MODE"); ok && val == "ON" {
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
	}

	uri := viper.GetString("DOCKER_REGISTRY_URI")
	if uri == "" {
		//panic("private registry uri not configured!")
		uri = "index.docker.io"
	}

	username := viper.GetString("DOCKER_REGISTRY_USERNAME")
	if username == "" {
		panic("private regsitry username not configured!")
	}

	password := viper.GetString("DOCKER_REGISTRY_PASSWORD")
	if password == "" {
		panic("private registry password not configured!")
	}

	PrivateRepoConfig = DockerRegistryConfig{
		URI:      uri,
		Username: username,
		Password: password,
	}
}
