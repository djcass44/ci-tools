# Provenance

Provenance is captured in the [SLSA](https://slsa.dev) format.

It supports the `0.2` and `1.0` formats.
The `1.0` format is preferred but isn't supported by many tools (e.g. Cosign only supports `0.2` as of Oct 2023).

## Contents

The provenance format attempts to follow the specification set by SLSA, however some of the fields are implementation-specific or up for interpretation.

### Inputs

The inputs captured are:

* The repository URL and commit SHA
* The image used by the CI runner and its digest
* The image used as the base layer and its digest
* The build configuration (e.g. `.gitlab-ci.yml`) file and its hash

### Commands

In order to accurately capture the build process, the provenance captures:

1. The arguments passed to the `ci` CLI applications
1. The arguments that `ci` passes to the underlying build tool.
1. The shell that was used by the CI runner (e.g. `/bin/bash`)

The combination of this data is all that is required to accurately capture how the application was built.

### Outputs

The primary output is the `subject` that tracks the built container image and its digest.

Additionally, `ci` also saves the `subject` value to a `build.txt` file in the working directory.
This file is useful for avoiding the round-trip to the container registry to get the digest of the container that was just built.

An example usage is to pass the image digest directly to Cosign.
