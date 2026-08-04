# docker

Module path: `github.com/global-torque/go-common/docker/v2`

Build seed for the shared go-common Docker image. This directory is not a
reusable go-common library package.

## Use For

- Building `cr.webdevelop.pro/global-torque/go-common` images.
- Pre-downloading heavy Go dependencies into a shared builder image.
- Precompiling `evm-api` dependencies into the Compose-mounted Go module and
  build cache paths.
- Shipping common `etc` files such as `make.sh`, `golangci.yml`, `air.toml`,
  and `pre-commit`.

## Do Not Use For

- Service package imports.
- Application runtime logic.

## Key Files

- `docker/Dockerfile`
- `docker/main.go`
- `docker/build-deploy.sh`
- `docker/evm-api-cache/evm-api.mod`
- `docker/evm-api-cache/evm-api.sum`
- `docker/evm-api-cache/dependencies.go`
- `docker/etc/make.sh`
- `docker/etc/golangci.yml`
- `docker/etc/air.toml`
- `docker/etc/pre-commit`

## Build Configuration

Docker build args:

- `GIT_COMMIT`
- `BUILD_DATE`
- `SERVICE_NAME`
- `REPOSITORY`
- `VERSION`
- `GOLANGCI_LINT_VERSION`
- `GOSEC_VERSION`
- `GCI_VERSION`

The image currently uses Go `1.25.8` on Alpine and installs `golangci-lint`,
`gosec`, and `gci`.

`docker/evm-api-cache/evm-api.mod` mirrors the consumer's module graph and
`dependencies.go` blank-imports its direct external packages behind the
`evmapi_cache` build tag. During the image build, Go downloads and compiles
that graph into `/go/pkg/mod` and
`/go/build-cache`. Docker/Podman initializes new named volumes from those
directories, so `evm-api` reuses the cache on its first Compose build.

## Wiring Pattern

Dependent service Dockerfiles can use the published image as a builder:

```Dockerfile
FROM cr.webdevelop.pro/global-torque/go-common:latest-dev AS builder
RUN ./make.sh build
```

Refresh the cache inputs after `evm-api/go.mod` or its direct external imports
change. The snapshot must match the consumer's selected versions exactly so
Go module and build cache keys remain reusable.

## CI

Root GitHub Actions run vet/tests with PostgreSQL and Pub/Sub emulator services,
then build and push the Docker image.

## Gotchas

- `docker/go.mod` module path is currently `github.com/global-torque/go-common/docker/v2`.
- `docker/main.go` blank-imports common dependencies to warm the image cache.
- Cache reuse depends on consumers using `/go/pkg/mod` and
  `/go/build-cache`; changing those paths bypasses the pre-warmed cache.
- This module is a build artifact, not a public package API.
