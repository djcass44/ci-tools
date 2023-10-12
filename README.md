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

Information about recipes can be found in [RECIPES.md](./docs/RECIPES.md).

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

Information about environment variables can be found in [ENVIRONMENT.md](./docs/ENVIRONMENT.md).

### Signing

Information about signing can be found in [COSIGN.md](./docs/COSIGN.md).

### Provenance

All execution of the build phase generate the following provenance:
* `sbom.cdx.json` - Software Bill of Materials (SBOM) in CycloneDX format
* `provenance.slsa.json` - SLSA build provenance in InToto format
* `build.txt` - text file containing the full path of the built image

These files are output to the build root and should be attested using Cosign.

Information about the build provenance can be found in [PROVENANCE.md](./docs/PROVENANCE.md).
