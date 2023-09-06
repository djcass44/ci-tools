package ociutil

import (
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	"github.com/chrismellard/docker-credential-acr-env/pkg/credhelper"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"log"
)

// BasicKeychain is an authn.Keychain implementation that
// uses basic auth to talk to a registry.
type BasicKeychain struct {
	auth Auth
}

type Auth struct {
	Registry string
	Username string
	Password string
}

func NewBasicKeychain(auth Auth) *BasicKeychain {
	return &BasicKeychain{
		auth: auth,
	}
}

func (b *BasicKeychain) Resolve(resource authn.Resource) (authn.Authenticator, error) {
	if resource.RegistryStr() != b.auth.Registry {
		log.Printf("registry does not match (expected: '%s', actual: '%s')", b.auth.Registry, resource.RegistryStr())
		return authn.Anonymous, nil
	}
	return &authn.Basic{
		Username: b.auth.Username,
		Password: b.auth.Password,
	}, nil
}

func KeyChain(auth Auth) authn.Keychain {
	keychains := []authn.Keychain{
		authn.DefaultKeychain,
		google.Keychain,
		authn.NewKeychainFromHelper(ecr.NewECRHelper(ecr.WithClientFactory(api.DefaultClientFactory{}))),
		authn.NewKeychainFromHelper(credhelper.NewACRCredentialsHelper()),
	}
	if auth.Username != "" && auth.Password != "" {
		keychains = append([]authn.Keychain{NewBasicKeychain(auth)}, keychains...)
	}
	return authn.NewMultiKeychain(keychains...)
}
