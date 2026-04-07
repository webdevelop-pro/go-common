FROM public.ecr.aws/docker/library/golang:1.25-alpine3.23

ARG GIT_COMMIT=unspecified
ARG BUILD_DATE=unspecified
ARG SERVICE_NAME=unspecified
ARG REPOSITORY=unspecified
ARG VERSION
ARG GOLANGCI_LINT_VERSION=v2.11.4

LABEL GIT_COMMIT=$GIT_COMMIT
LABEL BUILD_DATE=$BUILD_DATE
LABEL SERVICE_NAME=$SERVICE_NAME
LABEL REPOSITORY=$REPOSITORY

ENV GOPATH=/go
ENV VERSION=${VERSION:-unknown}

RUN apk update && apk upgrade --no-cache
RUN apk add --no-cache make gcc musl-dev linux-headers git gettext ca-certificates curl

ADD . /app

WORKDIR /app

RUN go mod download
RUN go build main.go
RUN wget -O - https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}

RUN go install github.com/securego/gosec/v2/cmd/gosec@latest
RUN go install github.com/daixiang0/gci@latest
RUN ln -s /usr/local/go/bin/go /usr/local/bin/go

COPY etc ./etc /app/etc/
CMD ["/bin/sh"]
