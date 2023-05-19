FROM ubuntu:20.04 as builder

ARG ARCH
ARG GO_VERSION
ARG ENV
ARG GO_PACKAGE=go$GO_VERSION.linux-$ARCH
ARG GOROOT=/usr/local/go

ENV GOPATH $GOROOT/packages
ENV PATH ${GOROOT}/bin:${GOPATH}/bin:${PATH}
ENV ENV=${ENV}

RUN apt-get update && apt-get install -y \
  ca-certificates \
  wget

# Setting up Go
RUN wget https://go.dev/dl/${GO_PACKAGE}.tar.gz 
RUN tar -C /usr/local -zxvf ${GO_PACKAGE}.tar.gz  

WORKDIR /app
COPY . .
RUN go build -o main main.go

COPY ${ENV} .
COPY start.sh .
COPY wait-for.sh .
COPY staging-ca-certificate.crt /etc/ssl/certs
COPY prod-ca-certificate.crt /etc/ssl/certs
COPY simple-bank-firebase-authentication.json .
COPY db/migration ./db/migration

EXPOSE 9090
CMD [ "/app/main" ]
