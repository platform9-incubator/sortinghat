#
# Python Dockerfile

# Pull base image.
FROM alpine:latest

# Update certificates
RUN apk --update upgrade && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

ADD ./bin/main /opt/pf9/whistle-ingest

# Define default command.
CMD ["/opt/pf9/whistle-ingest/bin/main"]
