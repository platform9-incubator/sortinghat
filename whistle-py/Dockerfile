#
# Python Dockerfile

# Pull base image.
FROM alpine:latest

# Update certificates
RUN apk --update upgrade && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

# Install Python.
RUN \
  apk update && \
  apk add python && \
  apk add py-setuptools && \
  mkdir -p /opt/pf9/whistle-log

ADD . /opt/pf9/whistle-log
RUN \
  cd /opt/pf9/whistle-log && \
  python setup.py install

# Define default command.
CMD ["/opt/pf9/whistle-log/web_server.py"]
