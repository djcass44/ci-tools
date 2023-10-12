# Environment

The `ci` application pulls almost all of its configuration from the environment.
The available options are separated into Ambient and Explicit environment variables.

## Explicit

Explicit environment variables are prefixed with `BUILD_` and are intended to provide the end user with a limited degree of control over the build process.

Explicit environment variables can be found in the [`v1/types.go`](../internal/api/v1/types.go) file.

## Ambient

Ambient environment variables are set by the CI runner and are used to get an understanding of the environment that `ci` is running in.

These variables are dependent on the CI tool but can be set or overriden manually if required.

Ambient environment variables can be found in the [`ctx`](../internal/api/ctx/) module.
