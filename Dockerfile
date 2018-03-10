FROM golang:1.8-alpine3.6

ENV TERRAFORM_VERSION=0.10.0

RUN apk add --update git bash openssh alpine-sdk

ENV TF_DEV=true
ENV TF_RELEASE=true

WORKDIR $GOPATH/src/github.com/hashicorp/terraform
RUN git clone https://github.com/hashicorp/terraform.git ./ && \
    git checkout v${TERRAFORM_VERSION} && \
    /bin/bash scripts/build.sh

ENV PROVIDER_NAME=terraform-provider-pagerduty
ENV PROVIDER_DIR=${GOPATH}/src/github.com/terraform-providers/${PROVIDER_NAME}
COPY . ${PROVIDER_DIR}

WORKDIR ${PROVIDER_DIR}
RUN make build

ENTRYPOINT ["sh"]
