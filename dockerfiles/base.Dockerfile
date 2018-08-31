FROM golang:1.11-alpine as go-builder

RUN apk --no-cache add \
    ca-certificates \
    git \
    wget

WORKDIR /app
COPY go.mod ./
COPY cmd/ ./cmd/

RUN go mod download

RUN CGO_ENABLED=0 go install \
    ./cmd/csv-splitter \
    ./cmd/robocopy \
    github.com/baltimore-sun-data/boreas

FROM alpine as go-curl
RUN apk --no-cache add \
    ca-certificates \
    curl

FROM go-curl as go-hugo
ARG HUGO_VERSION=0.48
RUN curl -L https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_Linux-64bit.tar.gz | tar xz -C /bin/

FROM go-curl as go-s3deploy
ARG S3DEPLOY_VERSION=2.0.1
RUN curl -L https://github.com/bep/s3deploy/releases/download/v${S3DEPLOY_VERSION}/s3deploy_${S3DEPLOY_VERSION}_Linux-64bit.tar.gz | tar xz -C /bin/

FROM ubuntu:18.04 as yarn-builder
RUN apt-get update -qq && \
    apt-get install -qq -y curl build-essential && \
    curl -sL https://deb.nodesource.com/setup_8.x | bash - && \
    apt-get install -qq -y nodejs && \
    npm install -g yarn

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn

COPY --from=go-builder /go/bin/boreas /bin/
COPY --from=go-builder /go/bin/csv-splitter /bin/
COPY --from=go-builder /go/bin/robocopy /bin/
COPY --from=go-hugo /bin/hugo /bin/
COPY --from=go-s3deploy /bin/s3deploy /bin/

COPY . .
RUN yarn run build
