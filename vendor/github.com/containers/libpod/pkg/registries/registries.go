package registries

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/containers/image/pkg/sysregistries"
	"github.com/containers/image/types"
	"github.com/containers/libpod/pkg/rootless"
	"github.com/docker/distribution/reference"
	"github.com/pkg/errors"
)

// userRegistriesFile is the path to the per user registry configuration file.
var userRegistriesFile = filepath.Join(os.Getenv("HOME"), ".config/containers/registries.conf")

// SystemRegistriesConfPath returns an appropriate value for types.SystemContext.SystemRegistriesConfPath
// (possibly "", which is not an error), taking into account rootless mode and environment variable overrides.
//
// FIXME: This should be centralized in a global SystemContext initializer inherited throughout the code,
// not haphazardly called throughout the way it is being called now.
func SystemRegistriesConfPath() string {
	if envOverride := os.Getenv("REGISTRIES_CONFIG_PATH"); len(envOverride) > 0 {
		return envOverride
	}

	if rootless.IsRootless() {
		if _, err := os.Stat(userRegistriesFile); err == nil {
			return userRegistriesFile
		}
	}

	return ""
}

// GetRegistries obtains the list of registries defined in the global registries file.
func GetRegistries() ([]string, error) {
	searchRegistries, err := sysregistries.GetRegistries(&types.SystemContext{SystemRegistriesConfPath: SystemRegistriesConfPath()})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse the registries.conf file")
	}
	return searchRegistries, nil
}

// GetInsecureRegistries obtains the list of insecure registries from the global registration file.
func GetInsecureRegistries() ([]string, error) {
	registries, err := sysregistries.GetInsecureRegistries(&types.SystemContext{SystemRegistriesConfPath: SystemRegistriesConfPath()})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse the registries.conf file")
	}
	return registries, nil
}

// GetRegistry returns the registry name from a string if specified
func GetRegistry(image string) (string, error) {
	// It is possible to only have the registry name in the format "myregistry/"
	// if so, just trim the "/" from the end and return the registry name
	if strings.HasSuffix(image, "/") {
		return strings.TrimSuffix(image, "/"), nil
	}
	imgRef, err := reference.Parse(image)
	if err != nil {
		return "", err
	}
	return reference.Domain(imgRef.(reference.Named)), nil
}
