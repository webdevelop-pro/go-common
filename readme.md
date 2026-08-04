# go-common

Shared, multi-module Go libraries and the development builder image used by
Global Torque services.

## Go Modules

There is no root Go module. Each top-level component has its own `go.mod` and
is published independently with a component-prefixed semantic version tag:

```text
github.com/global-torque/go-common/db/v2     -> db/v2.0.2
github.com/global-torque/go-common/orm/v2    -> orm/v2.0.2
github.com/global-torque/go-common/queue/v2  -> queue/v2.0.10
```

Use `go.work` for local cross-module development. Run repository validation
from the root:

```sh
./make.sh vet
./make.sh lint
./make.sh test
```

For a release, validate the exact component commit, create an annotated
`<component>/vX.Y.Z` tag, push that ref, and verify that the version resolves
with `GOWORK=off go mod download -json`.

## Development Builder Image

The image is published as:

```text
cr.webdevelop.pro/global-torque/go-common:latest-dev
```

It contains Go 1.25.8, native build packages, shared lint/security tools,
go-common dependencies, and pre-warmed module/build caches for `evm-api`.
The cache is stored at `/go/pkg/mod` and `/go/build-cache`, matching both the
`evm-api` builder and its Compose named-volume mount points.

Build and publish the default amd64/arm64 image:

```sh
./docker/build-deploy.sh
```

To publish only one platform (useful for a dev-only refresh):

```sh
PLATFORMS=linux/amd64 ./docker/build-deploy.sh
```

The script publishes both the short commit tag and `latest-dev`. The root CI
also builds the multi-platform image after vet and tests pass on `master`.

When `evm-api/go.mod` or its direct external imports change, refresh
`docker/evm-api-cache/evm-api.mod`, `docker/evm-api-cache/evm-api.sum`, and
`docker/evm-api-cache/dependencies.go` before publishing a new image.

## Using The Image

```Dockerfile
FROM cr.webdevelop.pro/global-torque/go-common:latest-dev AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ./make.sh build
```

See [docker/readme.md](docker/readme.md) for image details and `AGENTS.md` for
the package map and release safeguards.
