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
* [Ko](https://github.com/ko-build/ko) (`com.github.google.ko`)
* [Jib](https://github.com/GoogleContainerTools/jib) (`com.google.cloud.tools.jib-maven-plugin`)
* [BuildKit](https://github.com/moby/buildkit) (`com.github.moby.buildkit`)

### Usage:

```shell
ci build --recipe <recipe-name>
```

### Configuration

Command-line arguments:
* `--recipe` - name of the recipe to run
* `--recipe-template` - (*optional*) provide a custom recipe template file
* `--skip-docker-cfg` - (*optional*) disables the creation of the OCI registry credentials file, even if requested by a recipe
* `--skip-sbom` - (*optional*) disables the creation of the Software Bill of Materials (SBOM) file
* `--skip-slsa` - (*optional*) disables the creation of SLSA provenance
* `--skip-cosign-verify` - (*optional*) disables signature verification of the parent image
* `--cosign-verify-key` - (*optional*) path to the Cosign public key to use when verifying the parent image

Environment variables:
* `BUILD_IMAGE_PARENT` - the container image to use for the application runtime
* `BUILD_DOCKERFILE` - (*optional*) name of the `Dockerfile` within the build context
* `BUILD_ARG_*` - (*optional*) arbitrary key-value pairs to be passed to Dockerfile-based recipes
* `BUILD_CACHE_ENABLED` - (*optional*) enable or disable caching (default: `true`). Cache logic depends on the recipe
* `BUILD_CACHE_PATH` - (*optional*) the path that cache files will be stored (default `<project-root>/.cache`)
* `BUILD_GO_IMPORTPATH` - (*optional*) the import path used by Go projects. Useful when the `main.go` file is in a subdirectory
* `PROJECT_PATH` - (*optional*) path to the project within the build context (use for mono-repos)

### Provenance

All execution of the build phase generate the following provenance:
* `sbom.cdx.json` - Software Bill of Materials (SBOM) in CycloneDX format
* `provenance.slsa.json` - SLSA build provenance in InToto format

These files are output to the build root and should be attested using Cosign.
