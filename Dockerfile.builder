FROM debian:11

RUN apt update && apt install -y \
    gcc \
    gcc-x86-64-linux-gnu \
    gcc-aarch64-linux-gnu \
    git \
    make \
    upx \
    curl \
    && rm -rf /var/lib/apt/lists/*

RUN curl -L https://go.dev/dl/go1.26.3.linux-amd64.tar.gz | tar -C /usr/local -xz
ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /app