FROM public.ecr.aws/docker/library/golang:1.24.5-alpine3.21

ARG GIT_COMMIT=unspecified
ARG BUILD_DATE=unspecified
ARG SERVICE_NAME=unspecified
ARG REPOSITORY=unspecified
ARG VERSION

LABEL GIT_COMMIT=$GIT_COMMIT
LABEL BUILD_DATE=$BUILD_DATE
LABEL SERVICE_NAME=$SERVICE_NAME
LABEL REPOSITORY=$REPOSITORY

ENV GOPATH=/go
ENV VERSION=${VERSION:-unknown}

RUN apk add --no-cache make gcc musl-dev linux-headers git gettext ca-certificates curl

ADD . /app

WORKDIR /app

RUN go build main.go
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0

RUN go install github.com/securego/gosec/v2/cmd/gosec@latest
RUN go install github.com/daixiang0/gci@latest
RUN ./build-deploy.sh download
RUN ln -s /usr/local/go/bin/go /usr/local/bin/go

COPY etc ./etc /app/etc/
CMD ["/bin/sh"]
