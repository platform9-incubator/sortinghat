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
  apk add python py-pip && \
  apk add py-setuptools && \
  mkdir -p /opt/pf9/sortinghat-ui/logs

# Install gunicorn
RUN apk --update add py-gunicorn

ADD . /opt/pf9/sortinghat-ui
RUN \
  cd /opt/pf9/sortinghat-ui && \
  python setup.py install

COPY run.sh /opt/pf9/sortinghat-ui
WORKDIR /opt/pf9/sortinghat-ui

# Define default command.
CMD ["/opt/pf9/sortinghat-ui/run.sh"]
