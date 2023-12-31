################
# BUILD ZETAVISOR
################
FROM golang:1.20 AS cosmovisor
ENV GOBIN=/go/bin
ENV GOPATH=/go
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN apt update && \
    apt install unzip wget -y && \
    wget https://github.com/zeta-chain/cosmos-sdk/archive/refs/heads/zetavisor-v0.1.5.zip && \
    unzip zetavisor-v0.1.5.zip && \
    cd cosmos-sdk-zetavisor-v0.1.5/cosmovisor && \
    go get github.com/zeta-chain/cosmos-sdk/cosmovisor/cmd/zetavisor && \
    go install github.com/zeta-chain/cosmos-sdk/cosmovisor/cmd/zetavisor

FROM ubuntu:latest
RUN apt-get update && apt-get install -y curl ca-certificates
ENV GO_VERSION 1.20
RUN curl -L https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz -o go.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
EXPOSE 1317
EXPOSE 26656
EXPOSE 26657
COPY --from=cosmovisor /go/bin/zetavisor /usr/bin/zetavisor
COPY ./configs/* /app_version/
COPY ./scripts/network-genesis.sh /scripts/network-genesis.sh
COPY ./upgrades/ /app_version/

RUN apt-get update && \
    apt-get install libusb-dev build-essential gawk psmisc wget jq vim bash curl bison python3 python3-pip -y && \
    chmod -R a+x /scripts && \
    chmod -R 766 /scripts && \
    chmod -R a+x /app_version/ && \
    chmod -R 766 /app_version

CMD bash /scripts/network-genesis.sh