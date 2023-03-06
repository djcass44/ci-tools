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

Environment variables:
* `BUILD_IMAGE_PARENT` - the container image to use for the application runtime
* `BUILD_DOCKERFILE` - (*optional*) name of the `Dockerfile` within the build context
* `BUILD_ARG_*` - (*optional*) arbitrary key-value pairs to be passed to Dockerfile-based recipes
* `PROJECT_PATH` - (*optional*) path to the project within the build context (use for mono-repos)
