#
# Python Dockerfile

# Pull base image.
FROM alpine:3.1

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


COPY /tmp/whistle-log-wheel/  /opt/pf9/whistle-log

RUN \
  cd /opt/pf9/whistle-log && \
  pip  install *.whl

# Define default command.
CMD ["/opt/pf9/whistle-log/web_server.py"]
