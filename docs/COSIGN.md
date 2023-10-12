# Cosign

## Verification

Cosign verification can be performed automatically before the build is initiated.
It can be used to verify that a parent image has been signed.

> Note: parent image verification is useless with build tools that don't use the provided parent image (e.g., BuildKit).

Verifying the parent image is as simple as:

```shell
ci build --recipe com.github.google.ko --cosign-verify-key=/path/to/cosign.pub
```

### Verifying against multiple keys

When using `ci`, it's highly likely you will want to verify images that have been signed by different keys.
Rather than providing the `cosign.pub` file directly, you can instead provide a directory of keys:

```shell
ci build --recipe com.github.google.ko --cosign-verify-dir=/path/to/keys/
```

If at least 1 of the keys within the directory signed the parent image, `ci` will allow the build.

### Disabling verification

Verification can be disabled by setting the `--skip-cosign-verify` flag.

### Offline verification

Cosign verification is assumed to be offline.
If not, you can set the `--cosign-offline=false` flag.

When GitLab CI supports keyless signing, this may change to be the default.

## Signing

`ci` does not sign the artefacts that it generates.
Signing the artefacts is an opinionated process that is best left to the end user.

It does provide a `build.txt` that contains the full path to the artefact, which can be used as follows:

```shell
IMAGE="$(cat build.txt)"

echo "Signing image: '$IMAGE'"
# sign the container
cosign sign --key "$COSIGN_PRIVATE_KEY" "$IMAGE"
# sign the SLSA provenance and attach it to the image
cosign attest --key "$COSIGN_PRIVATE_KEY" --type slsaprovenance --predicate provenance.slsa.json "$IMAGE"
# sign the SBOM and attach it to the image
cosign attest --key "$COSIGN_PRIVATE_KEY" --type cyclonedx --predicate sbom.cdx.json "$IMAGE"
```
