#!/usr/bin/env bash
set -euo pipefail

readonly DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR" || exit 1

readonly IMAGE_NAME="bomfactory"
readonly ARCH="x86_64"

build-image() {
  local IMAGE_TAG="${IMAGE_TAG:-"latest"}"
  local ORG="bitbomdev"
  local MELANGE_IMAGE_REF="cgr.dev/chainguard/melange:latest"
  local APKO_IMAGE_REF="cgr.dev/chainguard/apko:latest"
  local SIGNING_KEY="melange.rsa.pub"
  local MELANGE_CONFIG="melange.yaml"
  local APKO_CONFIG="apko.yaml"
  local PACKAGES_DIR="packages-${IMAGE_NAME}"

  # Check if required files exist
  for file in "${MELANGE_CONFIG}" "${APKO_CONFIG}"; do
    if [[ ! -f "${file}" ]]; then
      echo "Error: ${file} not found in the current directory."
      exit 1
    fi
  done

  docker run --rm -v "${PWD}:/work" -w /work "${MELANGE_IMAGE_REF}" keygen

  docker run --rm --privileged -v "${PWD}:/work" -w /work \
    "${MELANGE_IMAGE_REF}" build ${MELANGE_CONFIG} \
    --arch "${ARCH}" \
    --source-dir /work \
    --out-dir "${PACKAGES_DIR}" \
    --signing-key melange.rsa

  # Build the image using APKO
  docker run --rm -v "${PWD}:/work" -w /work \
    --platform "linux/${ARCH}" \
    "${APKO_IMAGE_REF}" build "${APKO_CONFIG}" \
    "${ORG}/${IMAGE_NAME}:${IMAGE_TAG}" \
    "output-${IMAGE_NAME}-${ARCH}.tar" \
    --arch "${ARCH}" \
    -k "${SIGNING_KEY}" \
    -r "${PACKAGES_DIR}"

  docker load < "output-${IMAGE_NAME}-${ARCH}.tar"

  echo "Loaded image for architecture: ${ARCH}"

  # Tag the image
  docker tag "${ORG}/${IMAGE_NAME}:${IMAGE_TAG}"-amd64 "ghcr.io/${ORG}/${IMAGE_NAME}:latest"
}

build-image