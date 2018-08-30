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
    github.com/carlmjohnson/scattered/cmd/scattered

FROM alpine as go-curl
RUN apk --no-cache add \
    ca-certificates \
    curl

FROM go-curl as go-hugo
ARG HUGO_VERSION=0.48
RUN curl -L https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_Linux-64bit.tar.gz | tar xz -C /bin/

# Node comes with yarn
FROM node:9-alpine as yarn-builder
WORKDIR /go/src/github.com/baltimore-sun-data/voter-guide-2018
COPY package.json yarn.lock ./
RUN yarn

COPY --from=go-builder /go/bin/csv-splitter /bin/
COPY --from=go-builder /go/bin/robocopy /bin/
COPY --from=go-builder /go/bin/scattered /bin/
COPY --from=go-hugo /bin/hugo /bin/

COPY . .
RUN yarn run build

CMD yarn run test
