FROM golang:1.11-alpine as go-builder

RUN apk --no-cache add ca-certificates

WORKDIR /deploy

COPY --from=voter-guide-2018:base /bin/boreas /bin/
COPY --from=voter-guide-2018:base /bin/s3deploy /bin/
COPY --from=voter-guide-2018:base /app/public/ /deploy/public

# Add deploy script
COPY deploy.sh /bin/deploy.sh
RUN chmod +x /bin/deploy.sh

COPY s3deploy.yaml .

CMD /bin/deploy.sh
