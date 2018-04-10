FROM golang:1.10-alpine as go-builder

RUN apk --no-cache add \
    ca-certificates \
    git \
    wget

# Using wget for caching, see https://github.com/moby/moby/issues/15717
RUN wget -O /usr/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64
RUN chmod +x /usr/bin/dep

WORKDIR /go/src/github.com/baltimore-sun-data/voter-guide-2018
COPY Gopkg.toml Gopkg.lock ./
COPY cmd/ ./cmd/

RUN dep ensure -v

RUN go install \
    ./cmd/csv-splitter \
    ./vendor/github.com/carlmjohnson/scattered/cmd/scattered

FROM alpine as go-curl
RUN apk --no-cache add \
    ca-certificates \
    curl

FROM go-curl as go-hugo
ARG HUGO_VERSION=0.38.2
RUN curl -L https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_Linux-64bit.tar.gz | tar xz -C /bin/

FROM go-curl as go-minify
ARG MINIFY_VERSION=2.3.4
# Using 386 version because 64 bit version require .so missing in Alpine
RUN curl -L https://github.com/tdewolff/minify/releases/download/v${MINIFY_VERSION}/minify_${MINIFY_VERSION}_linux_386.tar.gz | tar xz -C /bin/

# Node comes with yarn
FROM node:9-alpine as yarn-builder
WORKDIR /go/src/github.com/baltimore-sun-data/voter-guide-2018
COPY package.json yarn.lock ./
RUN yarn

COPY --from=go-builder /go/bin/csv-splitter /bin/
COPY --from=go-builder /go/bin/scattered /bin/
COPY --from=go-hugo /bin/hugo /bin/
COPY --from=go-minify /bin/minify /bin/

COPY . .
RUN yarn run build

CMD yarn run test
