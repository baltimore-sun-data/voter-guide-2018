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
    ./vendor/github.com/carlmjohnson/scattered/cmd/scattered \
    ./vendor/github.com/baltimore-sun-data/boreas \
    ./vendor/github.com/bep/s3deploy \
    ./vendor/github.com/spf13/hugo \
    ./vendor/github.com/tdewolff/minify/cmd/minify

# Node comes with yarn
FROM node:9-alpine as yarn-builder
WORKDIR /go/src/github.com/baltimore-sun-data/voter-guide-2018
COPY package.json yarn.lock ./
RUN yarn

COPY --from=go-builder /go/bin/csv-splitter /bin/
COPY --from=go-builder /go/bin/scattered /bin/
COPY --from=go-builder /go/bin/hugo /bin/
COPY --from=go-builder /go/bin/minify /bin/

COPY . .
RUN yarn run build

# Create final container
FROM go-builder as deployer

WORKDIR /deploy

COPY --from=go-builder /go/bin/boreas /bin/
COPY --from=go-builder /go/bin/s3deploy /bin/
COPY --from=yarn-builder /go/src/github.com/baltimore-sun-data/voter-guide-2018/public/ /deploy/public

# Add deploy script
COPY deploy.sh /bin/deploy.sh
RUN chmod +x /bin/deploy.sh

COPY s3deploy.yaml .

CMD /bin/deploy.sh
