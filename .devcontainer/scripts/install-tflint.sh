#!/usr/bin/env bash
# Install tflint without cosign/Rekor verification (the terraform devcontainer
# feature uses keyless cosign which depends on rekor.sigstore.dev availability).
# Verification is done via the SHA256 checksums file published with each release.
set -euo pipefail

TFLINT_VERSION="${TFLINT_VERSION:-0.62.1}"

case "$(uname -m)" in
	x86_64)            arch="amd64" ;;
	aarch64 | arm64)   arch="arm64" ;;
	*) echo "unsupported arch: $(uname -m)" >&2; exit 1 ;;
esac

tmp="$(mktemp -d)"
trap 'rm -rf "${tmp}"' EXIT

zip="tflint_linux_${arch}.zip"
base="https://github.com/terraform-linters/tflint/releases/download/v${TFLINT_VERSION}"

curl -fsSL -o "${tmp}/${zip}"            "${base}/${zip}"
curl -fsSL -o "${tmp}/checksums.txt"     "${base}/checksums.txt"

(cd "${tmp}" && sha256sum --ignore-missing -c checksums.txt)

sudo unzip -o "${tmp}/${zip}" -d /usr/local/bin
sudo chmod +x /usr/local/bin/tflint
tflint --version
