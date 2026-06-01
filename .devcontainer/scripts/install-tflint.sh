#!/usr/bin/env bash
# Install tflint without cosign/Rekor verification (the terraform devcontainer
# feature uses keyless cosign which depends on rekor.sigstore.dev availability).
# Verification is done against SHA256 checksums pinned in this script so the
# integrity check does not share a trust root with the downloaded artifact.
set -euo pipefail

# Pinned tflint version. Bump this value (and the matching checksums below)
# together when upgrading; the version is intentionally not overridable.
readonly TFLINT_VERSION="0.62.1"

# SHA256 checksums for the pinned version, copied from the upstream
# checksums.txt published alongside the release.
readonly SHA256_LINUX_AMD64="c004ec45ade3caf87cd4089feb1d2af9f7df57b13140a36df8a63c0a8cc69f14"
readonly SHA256_LINUX_ARM64="9f3bce43f7f58f05ddcf193f0bf1f7e7a9c7a79d7f46f72dd38e97f96fcaf14c"

case "$(uname -m)" in
	x86_64)            arch="amd64"; expected_sha="${SHA256_LINUX_AMD64}" ;;
	aarch64 | arm64)   arch="arm64"; expected_sha="${SHA256_LINUX_ARM64}" ;;
	*) echo "unsupported arch: $(uname -m)" >&2; exit 1 ;;
esac

tmp="$(mktemp -d)"
trap 'rm -rf "${tmp}"' EXIT

zip="tflint_linux_${arch}.zip"
base="https://github.com/terraform-linters/tflint/releases/download/v${TFLINT_VERSION}"

curl --proto '=https' --tlsv1.2 -fsSL -o "${tmp}/${zip}" "${base}/${zip}"

echo "${expected_sha}  ${zip}" > "${tmp}/checksums.txt"
(cd "${tmp}" && sha256sum -c checksums.txt)

sudo unzip -o "${tmp}/${zip}" -d /usr/local/bin
sudo chmod +x /usr/local/bin/tflint
tflint --version
