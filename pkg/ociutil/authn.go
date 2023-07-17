package ociutil

import (
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	"github.com/chrismellard/docker-credential-acr-env/pkg/credhelper"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/google"
)

// BasicKeychain is an authn.Keychain implementation that
// uses basic auth to talk to a registry.
type BasicKeychain struct {
	registry string
	username string
	password string
}

func NewBasicKeychain(registry, username, password string) *BasicKeychain {
	return &BasicKeychain{
		registry: registry,
		username: username,
		password: password,
	}
}

func (b *BasicKeychain) Resolve(resource authn.Resource) (authn.Authenticator, error) {
	if resource.RegistryStr() != b.registry {
		return authn.Anonymous, nil
	}
	return &authn.Basic{
		Username: b.username,
		Password: b.password,
	}, nil
}

func KeyChain(registry, username, password string) authn.Keychain {
	keychains := []authn.Keychain{
		authn.DefaultKeychain,
		google.Keychain,
		authn.NewKeychainFromHelper(ecr.NewECRHelper(ecr.WithClientFactory(api.DefaultClientFactory{}))),
		authn.NewKeychainFromHelper(credhelper.NewACRCredentialsHelper()),
	}
	if username != "" && password != "" {
		keychains = append([]authn.Keychain{NewBasicKeychain(registry, username, password)}, keychains...)
	}
	return authn.NewMultiKeychain(keychains...)
}
