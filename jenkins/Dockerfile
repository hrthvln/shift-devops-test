FROM jenkins/jenkins:lts-jdk21

USER root

RUN apt-get update \
    && apt-get install -y \
        docker.io \
        git \
        curl \
        openssh-client \
    && rm -rf /var/lib/apt/lists/*

COPY --from=golang:1.27.1 /usr/local/go /usr/local/go

ENV PATH="/usr/local/go/bin:${PATH}"

USER root