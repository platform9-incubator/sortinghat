#
# Python Dockerfile

# Pull base image.\
FROM blang/golang-alpine:latest

# Update certificates
RUN apk --update upgrade && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

RUN mkdir -p /opt/pf9
ADD ./whistle-log /opt/pf9/

# Define default command.
CMD ["/opt/pf9/whistle-log"]
