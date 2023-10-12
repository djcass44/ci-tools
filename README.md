# CI Tools

CI Tools is a CLI application designed for normalising tools commonly used in CI/CD pipelines.
Currently, the following phases are supported:
* `build`
* `test` (*planned*)

## Installing

Installation instructions can be found in the releases page.

## Build

The build phase supports a number of tools.
Tools are configured using "recipes", which are detailed in the [`recipes.tpl.yaml`](internal/api/v1/recipes.tpl.yaml) file.

The default recipe list contains:
* [Ko](https://github.com/ko-build/ko) (`com.github.google.ko`) for Go
* [Jib](https://github.com/GoogleContainerTools/jib) (`com.google.cloud.tools.jib-maven-plugin`) for Java
* [BuildKit](https://github.com/moby/buildkit) (`com.github.moby.buildkit`) for Dockerfile
* [Nib](https://github.com/djcass44/nib) (`com.github.djcass44.nib`) for Static web applications
* [`all-your-base`](https://github.com/djcass44/all-your-base) (`com.github.djcass44.all-your-base`) for base images

### Usage:

```shell
ci build --recipe <recipe-name>
```

While `ci` can be run manually, it expects to be run in a CI system and pulls information from the ambient environment variables provided by the CI system.

It currently supports:

* GitLab CI

### Configuration

#### Command-line arguments:

Add `--help` to any command to view the full set of options.

#### Environment variables:
* `BUILD_IMAGE_PARENT` - the container image to use for the application runtime
* `BUILD_DOCKERFILE` - (*optional*) name of the `Dockerfile` within the build context
* `BUILD_ARG_*` - (*optional*) arbitrary key-value pairs to be passed to Dockerfile-based recipes
* `BUILD_CACHE_ENABLED` - (*optional*) enable or disable caching (default: `true`). Cache logic depends on the recipe
* `BUILD_CACHE_PATH` - (*optional*) the path that cache files will be stored (default `<project-root>/.cache`)
* `BUILD_GO_IMPORTPATH` - (*optional*) the import path used by Go projects. Useful when the `main.go` file is in a subdirectory
* `PROJECT_PATH` - (*optional*) path to the project within the build context (useful for mono-repos)

### Provenance

All execution of the build phase generate the following provenance:
* `sbom.cdx.json` - Software Bill of Materials (SBOM) in CycloneDX format
* `provenance.slsa.json` - SLSA build provenance in InToto format
* `build.txt` - text file containing the full path of the built image

These files are output to the build root and should be attested using Cosign.

#### SLSA

Provenance is captured in the [SLSA](https://slsa.dev/) format.
It supports the `0.2` and `1.0` formats. The `1.0` format is recommended but doesn't support as many tools (e.g. Cosign only supports `0.2` as of 2023).
