# CI Tools

CI Tools is a CLI application designed for normalising tools commonly used in CI/CD pipelines.
Currently, the following phases are supported:
* `build`
* `test` (*planned*)

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
* `--recipe-template` - provide a custom recipe template file
* `--skip-docker-cfg` - disables the creation of the OCI registry credentials file, even if requested by a recipe
* `--skip-sbom` - disables the creation of the Software Bill of Materials (SBOM) file

Environment variables:
* `BUILD_IMAGE_PARENT` - the container image to use for the application runtime
* `BUILD_DOCKERFILE` - name of the `Dockerfile` within the build context
* `BUILD_ARG_*` - arbitrary key-value pairs to be passed to Dockerfile-based recipes
* `PROJECT_PATH` - path to the project within the build context (use for mono-repos)
