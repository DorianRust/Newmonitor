FROM ubuntu:16.04
FROM golang:1.11 as builder

ENV DEBIAN_FRONTEND noninteractive
ENV LANG C.UTF-8

# Default versions
ENV INFLUXDB_VERSION 1.6.4
ENV GRAFANA_VERSION  5.3.1

# Database Defaults
ENV INFLUXDB_GRAFANA_DB datasource
ENV INFLUXDB_GRAFANA_USER datasource
ENV INFLUXDB_GRAFANA_PW datasource
ENV GOPATH /root/go
ENV GOROOT /usr/local/go

ENV GF_DATABASE_TYPE=sqlite3

# Fix bad proxy issue
COPY docker/system/99fixbadproxy /etc/apt/apt.conf.d/99fixbadproxy

# Clear previous sources
RUN rm /var/lib/apt/lists/* -vf

# Base dependencies
RUN apt-get -y update && \
 apt-get -y dist-upgrade && \
 apt-get -y --force-yes install \
  apt-utils \
  ca-certificates \
  curl \
  git \
  htop \
  libfontconfig \
  nano \
  net-tools \
  openssh-server \
  supervisor \
  wget && \
 curl -sL https://deb.nodesource.com/setup_7.x | bash - && \
 apt-get install -y nodejs

WORKDIR /root

RUN mkdir -p /var/log/supervisor && \
    mkdir -p /var/run/sshd && \
    mkdir -p /root/go/src/github.com/sasaxie/monitor && \
    mkdir -p /root/go/bin/conf && \
    sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
    echo 'root:root' | chpasswd && \
    rm -rf .ssh && \
    rm -rf .profile && \
    mkdir .ssh

# Install InfluxDB
RUN wget https://dl.influxdata.com/influxdb/releases/influxdb_${INFLUXDB_VERSION}_amd64.deb && \
    dpkg -i influxdb_${INFLUXDB_VERSION}_amd64.deb && rm influxdb_${INFLUXDB_VERSION}_amd64.deb

# Install Grafana
RUN wget https://s3-us-west-2.amazonaws.com/grafana-releases/release/grafana_${GRAFANA_VERSION}_amd64.deb && \
    dpkg -i grafana_${GRAFANA_VERSION}_amd64.deb && rm grafana_${GRAFANA_VERSION}_amd64.deb

# Install Monitor
COPY . go/src/github.com/sasaxie/monitor

# Cleanup
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Configure Supervisord, SSH and base env
COPY docker/supervisord/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY docker/ssh/id_rsa .ssh/id_rsa
COPY docker/bash/profile .profile

# Configure InfluxDB
COPY docker/influxdb/influxdb.conf /etc/influxdb/influxdb.conf
COPY docker/influxdb/init.sh /etc/init.d/influxdb

# Configure Grafana
COPY docker/grafana/grafana.ini /etc/grafana/grafana.ini

RUN chmod 0755 /etc/init.d/influxdb

#Configure Monitor
RUN go install github.com/sasaxie/monitor
COPY conf/monitor.toml go/bin/conf
COPY conf/nodes.json go/bin/conf

CMD ["/usr/bin/supervisord"]