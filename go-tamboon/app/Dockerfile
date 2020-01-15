FROM golang:1.11.4

WORKDIR /go/src/go-tamboon
ENV OMISE_PUBLIC_KEY=pkey_test_521w1g1t7w4x4rd22z0
ENV OMISE_SECRET_KEY=skey_test_521w1g1t6yh7sx4pu8n
ENV GO111MODULE=on

COPY src/go.mod ./
COPY src/go.sum ./
RUN go mod download
