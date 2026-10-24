FROM debian:11

# Debian 11 已停止维护,默认源删了旧包,改用归档源
RUN sed -i 's|deb.debian.org/debian-security|archive.debian.org/debian-security|; s|deb.debian.org/debian|archive.debian.org/debian|; /bullseye-updates/d' /etc/apt/sources.list \
    && echo 'Acquire::Check-Valid-Until false;' > /etc/apt/apt.conf.d/99archive

RUN apt update && apt install -y \
    gcc \
    gcc-x86-64-linux-gnu \
    gcc-aarch64-linux-gnu \
    git \
    make \
    upx \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Go 工具链按宿主架构下载(arm64 Mac 原生运行,Intel 机器用 amd64)
RUN curl -L https://go.dev/dl/go1.26.3.linux-$(dpkg --print-architecture).tar.gz | tar -C /usr/local -xz
ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /app
